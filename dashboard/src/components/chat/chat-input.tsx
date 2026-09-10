"use client";

import { useState, type FormEvent } from "react";

type ChatInputProps = {
    onSend: (query: string) => void;
    disabled?: boolean;
};

export function ChatInput({ onSend, disabled }: ChatInputProps) {
    const [value, setValue] = useState("");

    function handleSubmit(e: FormEvent) {
        e.preventDefault();
        const trimmed = value.trim();
        if (!trimmed || disabled) return;

        onSend(trimmed);
        setValue("");
    }

    return (
        <form onSubmit={handleSubmit} className="flex gap-2 border-t border-gray-200 bg-white p-3">
            <input
                type="text"
                value={value}
                onChange={(e) => setValue(e.target.value)}
                disabled={disabled}
                placeholder="Tanyakan sesuatu..."
                className="flex-1 rounded-[10px] border border-gray-300 bg-white px-3 py-2 text-sm text-blue-600 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-600 disabled:opacity-50"
            />
            <button
                type="submit"
                disabled={disabled}
                className="rounded-md bg-blue-600 px-4 py-2 text-sm text-white disabled:opacity-50"
            >
                Kirim
            </button>
        </form>
    );
}