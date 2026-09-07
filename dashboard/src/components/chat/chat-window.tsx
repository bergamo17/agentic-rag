"use client";

import { useState } from "react";
import { useWidgetStore } from "@/lib/store/widgets-store";
import { sendChatMessage } from "@/lib/api-client";
import { MessageBubble} from "./message-bubble";
import { ChatInput } from "./chat-input";

type Message = {
    id: string;
    role: "user" | "agent";
    content: string;
    isPartial?: boolean;
}

export function ChatWindow() {
    const [messages, setMessages] = useState<Message[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const addWidgets = useWidgetStore((state) => state.addWidgets);

    async function handleSend(query: string) {
        const userMessage: Message = {
            id: crypto.randomUUID(),
            role: "user",
            content: query,
        };
        setMessages((prev) => [...prev, userMessage]);
        setIsLoading(true);
        setError(null);

        try {
            const result = await sendChatMessage(query);

            const agentMessage: Message = {
                id: crypto.randomUUID(),
                role: "agent",
                content: result.answer,
                isPartial: result.isPartial,
            };
            setMessages((prev) => [...prev, agentMessage]);

            if (result.widgets.length > 0) {
                addWidgets(result.widgets);
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : "Terjadi Kesalahan");
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <div className="flex flex-col h-full">
        <div className="flex-1 overflow-y-auto p-4 space-y-3">
            {messages.map((msg) => (
            <MessageBubble key={msg.id} message={msg} />
            ))}
            {isLoading && <p className="text-sm text-gray-400">Agent sedang berpikir...</p>}
            {error && <p className="text-sm text-red-500">Error: {error}</p>}
        </div>
        <ChatInput onSend={handleSend} disabled={isLoading} />
        </div>
    );
}