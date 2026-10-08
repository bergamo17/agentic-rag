"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useConversationStore } from "@/lib/store/conversations-store";
import { useWidgetStore } from "@/lib/store/widgets-store";
import { sendChatMessage, getConversationMessage, ApiError, type ChatMessage } from "@/lib/api-client";
import { MessageBubble} from "./message-bubble";
import { ChatInput } from "./chat-input";
import type { Widget } from "@/lib/widget-schema";
import type { GeneratedDocument, Message } from "@/lib/chat-type";

// type GeneratedDocument = {
//     title: string;
//     theme: string;
//     outputPath: string;
// }

// type Message = {
//     id: string;
//     role: "user" | "agent";
//     content: string;
//     attachment?: string[];
//     isPartial?: boolean;
//     widgets?: Widget[];
//     documents?: GeneratedDocument[];
// }

function toMessage(m: ChatMessage): Message {
    return {
        id: m.id,
        role: m.role === "assistant" ? "agent" : "user",
        content: m.content,
        isPartial: m.isPartial,
        widgets: m.widgets,
        documents: m.documents,
    };
}

export function ChatWindow() {
    const router = useRouter();
    const urlConversationId = useSearchParams().get("c");

    const [messages, setMessages] = useState<Message[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isHistoryLoading, setIsHistoryLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [warning, setWarning] = useState<string | null>(null);

    // conversation yang isinya sedang tampil di layar
    const syncedIdRef = useRef<string | null>(null);
    // naik setiap pindah conversation, untuk mengabaikan respons yang datang terlambat
    const sessionRef = useRef(0);

    const addWidgets = useWidgetStore((s) => s.addWidgets);
    const clearWidgets = useWidgetStore((s) => s.clearWidgets);
    const upsertConversation = useConversationStore((s) => s.upsert);

    useEffect(() => {
        // URL sudah sinkron dengan layar (mis. baru dibuat dari pesan pertama)
        if (urlConversationId === syncedIdRef.current) return;

        sessionRef.current += 1;
        syncedIdRef.current = null;
        setMessages([]);
        setError(null);
        setIsLoading(false);
        clearWidgets();

        if (!urlConversationId) return; // chat baru

        let cancelled = false;
        setIsHistoryLoading(true);

        getConversationMessage(urlConversationId)
            .then((history) => {
                if (cancelled) return;
                syncedIdRef.current = urlConversationId;
                setMessages(history.map(toMessage));

                const widgets = history.flatMap((m) => m.widgets);
                if (widgets.length > 0) addWidgets(widgets);
            })
            .catch((err) => {
                if (cancelled) return;
                if (err instanceof ApiError && err.status === 404) {
                    router.replace("/main/chat"); // conversation sudah dihapus
                    return;
                }
                setError(err instanceof Error ? err.message : "Gagal memuat riwayat chat");
            })
            .finally(() => {
                if (!cancelled) setIsHistoryLoading(false);
            });

        return () => {
            cancelled = true;
        };
    }, [urlConversationId, addWidgets, clearWidgets, router]);

    async function handleSend(query: string, files: File[]) {
        const session = sessionRef.current;

        setMessages((prev) => [
            ...prev,
            { id: crypto.randomUUID(), role: "user", content: query, attachment: files.map((f) => f.name) },
        ]);
        setIsLoading(true);
        setError(null);
        setWarning(null);

        try {
            const result = await sendChatMessage(query, syncedIdRef.current, files);

            // tetap tersimpan di backend, tapi jangan tampilkan di conversation lain
            upsertConversation({ id: result.conversationId, title: result.title });
            if (session !== sessionRef.current) return;

            if (result.skippedFiles.length > 0) {
                setWarning(`File is not processed: ${result.skippedFiles.join(", ")}`);
            }

            setMessages((prev) => [
                ...prev,
                {
                    id: crypto.randomUUID(),
                    role: "agent",
                    content: result.answer,
                    isPartial: result.isPartial,
                    widgets: result.widgets,
                    documents: result.documents,
                },
            ]);

            if (result.widgets.length > 0) addWidgets(result.widgets);

            if (syncedIdRef.current !== result.conversationId) {
                syncedIdRef.current = result.conversationId;
                router.replace(`/main/chat?c=${result.conversationId}`, { scroll: false });
            }
        } catch (err) {
            if (session === sessionRef.current) {
                setError(err instanceof Error ? err.message : "Terjadi Kesalahan");
            }
        } finally {
            if (session === sessionRef.current) setIsLoading(false);
        }
    }

    return (
        <div className="flex flex-col h-full">
            <div className="flex-1 overflow-y-auto p-4 space-y-3">
                {isHistoryLoading && <p className="text-sm text-gray-400">Memuat riwayat...</p>}
                {messages.map((msg) => (
                    <MessageBubble key={msg.id} message={msg} />
                ))}
                {isLoading && <p className="text-sm text-gray-400">Agent sedang berpikir...</p>}
                {error && <p className="text-sm text-red-500">Error: {error}</p>}
                {warning && <p className="text-sm text-amber-600">⚠ {warning}</p>}
            </div>
            <ChatInput onSend={handleSend} disabled={isLoading || isHistoryLoading} />
        </div>
    );
}