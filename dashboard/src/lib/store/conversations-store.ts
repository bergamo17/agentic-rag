import { create } from "zustand";
import {
    Conversation,
    listConversations,
    renameConversation,
    deleteConversation
} from "@/lib/api-client";

type ConversationState = {
    conversations: Conversation[];
    load: () => Promise<void>;
    upsert: (c: { id: string; title: string }) => void;
    rename: (id: string, title: string) => Promise<void>;
    remove: (id: string) => Promise<void>;
}

export const useConversationStore = create<ConversationState>((set) => ({
    conversations: [],

    load: async () => {
        try {
            set({ conversations: await listConversations() });
        } catch (err) {
            console.error("Gagal memuat daftar conversation", err);
        }
    },

    // dipanggil setelah pesan terkirim: conversation baru muncul / naik ke atas
    upsert: ({ id, title }) =>
        set((state) => {
            const existing = state.conversations.find((c) => c.id === id);
            const now = new Date().toISOString();
            const updated: Conversation = {
                id,
                title,
                createdAt: existing?.createdAt ?? now,
                updatedAt: now,
            };
            return {
                conversations: [updated, ...state.conversations.filter((c) => c.id !== id)],
            };
        }),

    rename: async (id, title) => {
        const updated = await renameConversation(id, title);
        set((state) => ({
            conversations: state.conversations.map((c) => (c.id === id ? updated : c)),
        }));
    },

    remove: async (id) => {
        await deleteConversation(id);
        set((state) => ({
            conversations: state.conversations.filter((c) => c.id !== id),
        }));
    },
}));