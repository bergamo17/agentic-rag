"use client";

import { useRef, useState, type ChangeEvent, type FormEvent } from "react";
import { Paperclip, X } from "lucide-react";

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
        setFiles([]);
    }

        return (
        <form onSubmit={handleSubmit} className="border-t border-border bg-card p-3">
            {files.length > 0 && (
                <div className="mb-2 flex flex-wrap gap-2">
                    {files.map((file, i) => (
                        <span
                            key={`${file.name}-${i}`}
                            className="flex items-center gap-1.5 rounded-tight-sm bg-primary-soft px-2 py-1 text-xs font-medium text-primary-accent"
                        >
                            {file.name}
                            <button
                                type="button"
                                onClick={() => removeFile(i)}
                                className="text-primary-accent/70 hover:text-danger"
                                aria-label={`Hapus ${file.name}`}
                            >
                                <X size={12} />
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
                    className="rounded-md border border-grey-300 px-3 py-2 text-grey-600 hover:bg-grey-100 disabled:opacity-50"
                    aria-label="Lampirkan file"
                >
                    <Paperclip size={16} />
                </button>

                <input
                    type="text"
                    value={value}
                    onChange={(e) => setValue(e.target.value)}
                    disabled={disabled}
                    placeholder="Tanyakan sesuatu..."
                    className="flex-1 rounded-tight-sm border border-grey-300 bg-card px-3 py-2 text-sm font-medium text-foreground placeholder:font-normal placeholder:text-grey-500 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/30 disabled:opacity-50"
                />
                <button
                    type="submit"
                    disabled={disabled}
                    className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-white shadow-[inset_0_1px_0_rgba(255,255,255,0.14),0_1px_2px_rgba(4,104,204,0.16),0_4px_10px_rgba(4,104,204,0.14)] hover:bg-primary-active disabled:opacity-50"
                >
                    Kirim
                </button>
            </div>
        </form>
    );
}