"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { Pencil, Trash2, Plus } from "lucide-react";
import { useConversationStore } from "@/lib/store/conversations-store";

export function ConversationList() {
    const router = useRouter();
    const activeId = useSearchParams().get("c");

    const { conversations, load, rename, remove } = useConversationStore();
    const [editingId, setEditingId] = useState<string | null>(null);
    const [draftTitle, setDraftTitle] = useState("");

    useEffect(() => {
        load();
    }, [load]);

    async function commitRename(id: string) {
        const title = draftTitle.trim();
        setEditingId(null);
        if (!title) return;
        try {
            await rename(id, title);
        } catch {
            alert("Gagal mengganti judul");
        }
    }

    async function handleDelete(id: string, title: string) {
        if (!confirm(`Hapus "${title}"?`)) return;
        try {
            await remove(id);
            if (id === activeId) router.replace("/main/chat");
        } catch {
            alert("Gagal menghapus percakapan");
        }
    }

        return (
        <div className="flex flex-col gap-1">
            <Link
                href="/main/chat"
                className="mb-2 flex items-center gap-2 rounded-md border border-white/30 px-3 py-2 text-sm font-medium text-white hover:bg-white/10"
            >
                <Plus size={16} className="shrink-0" />
                <span className="whitespace-nowrap">Chat baru</span>
            </Link>

            <p className="px-3 pb-1 text-xs font-medium text-[#8FC1EA]">Riwayat</p>

            {conversations.map((c) => {
                const isActive = c.id === activeId;

                if (editingId === c.id) {
                    return (
                        <input
                            key={c.id}
                            autoFocus
                            value={draftTitle}
                            maxLength={100}
                            onChange={(e) => setDraftTitle(e.target.value)}
                            onBlur={() => commitRename(c.id)}
                            onKeyDown={(e) => {
                                if (e.key === "Enter") commitRename(c.id);
                                if (e.key === "Escape") setEditingId(null);
                            }}
                            className="rounded-md border border-accent bg-white px-3 py-2 text-sm text-grey-900 outline-none"
                        />
                    );
                }

                return (
                    <div
                        key={c.id}
                        className={`group flex items-center rounded-md text-sm transition-colors ${
                            isActive ? "bg-white/15 text-white" : "text-[#C7E1F7] hover:bg-white/10"
                        }`}
                    >
                        <Link
                            href={`/main/chat?c=${c.id}`}
                            title={c.title}
                            className="min-w-0 flex-1 truncate px-3 py-2"
                        >
                            {c.title}
                        </Link>
                        <div className="hidden shrink-0 gap-1 pr-2 group-hover:flex">
                            <button
                                aria-label="Ganti nama"
                                onClick={() => {
                                    setEditingId(c.id);
                                    setDraftTitle(c.title);
                                }}
                                className="rounded p-1 hover:bg-white/15"
                            >
                                <Pencil size={13} />
                            </button>
                            <button
                                aria-label="Hapus"
                                onClick={() => handleDelete(c.id, c.title)}
                                className="rounded p-1 hover:bg-danger hover:text-white"
                            >
                                <Trash2 size={13} />
                            </button>
                        </div>
                    </div>
                );
            })}
        </div>
    );
}