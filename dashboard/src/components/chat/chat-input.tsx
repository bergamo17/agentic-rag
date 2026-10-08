"use client";

import { useRef, useState, type ChangeEvent, type FormEvent } from "react";

type ChatInputProps = {
    onSend: (query: string, files: File[]) => void;
    disabled?: boolean;
};

const MAX_FILES = 5;
const MAX_SIZE_MB = 10;

export function ChatInput({ onSend, disabled }: ChatInputProps) {
    const [value, setValue] = useState("");
    const [files, setFiles] = useState<File[]>([]);
    const fileInputRef = useRef<HTMLInputElement>(null);

    function handleFileChange(e: ChangeEvent<HTMLInputElement>) {
        const picked = Array.from(e.target.files ?? []);
        const valid = picked.filter((f) => f.size <= MAX_SIZE_MB * 1024 * 1024);

        setFiles((prev) => [...prev, ...valid].slice(0, MAX_FILES));

        // reset agar file yang sama bisa dipilih ulang
        e.target.value = "";
    }

    function removeFile(index: number) {
        setFiles((prev) => prev.filter((_, i) => i !== index));
    }

    function handleSubmit(e: FormEvent) {
        e.preventDefault();
        const trimmed = value.trim();
        if ((!trimmed && files.length === 0) || disabled) return;

        onSend(trimmed, files);
        setValue("");
    }

    return (
        <form onSubmit={handleSubmit} className="border-t border-gray-200 bg-white p-3">
            {files.length > 0 && (
                <div className="mb-2 flex flex-wrap gap-2">
                    {files.map((file, i) => (
                        <span
                            key={`${file.name}-${i}`}
                            className="flex items-center gap-1 rounded-md bg-gray-100 px-2 py-1 text-xs text-gray-700"
                        >
                            {file.name}
                            <button
                                type="button"
                                onClick={() => removeFile(i)}
                                className="text-gray-500 hover:text-red-600"
                                aria-label={`Hapus ${file.name}`}
                            >
                                ✕
                            </button>
                        </span>
                    ))}
                </div>
            )}

            <div className="flex gap-2">
                <input
                    ref={fileInputRef}
                    type="file"
                    multiple
                    hidden
                    accept="image/*,.pdf,.docx,.txt,.md"
                    onChange={handleFileChange}
                />
                <button
                    type="button"
                    onClick={() => fileInputRef.current?.click()}
                    disabled={disabled || files.length >= MAX_FILES}
                    className="rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-600 disabled:opacity-50"
                    aria-label="Lampirkan file"
                >
                    📎
                </button>

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
            </div>
        </form>
    );
}