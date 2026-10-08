import type { Widget } from "@/lib/widget-schema";

export type GeneratedDocument = {
    title: string;
    theme: string;
    outputPath: string;
    format: string;
};

export type Message = {
    id: string;
    role: "user" | "agent";
    content: string;
    attachments?: string[];
    isPartial?: boolean;
    widgets?: Widget[];
    documents?: GeneratedDocument[];
};