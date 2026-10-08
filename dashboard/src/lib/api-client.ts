import { parseWidget, Widget } from "@/lib/widget-schema";
import type { GeneratedDocument } from "@/lib/chat-type";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

type RawWidget = {
    id: string;
    widget_type: string;
    title: string;
    data: string;
};

type RawDocument = {
    title: string;
    theme: string;
    output_path: string;
    format: string;
};

type ChatAgentResponse = {
    conversation_id: string;
    title: string;
    user_message_id: string;
    message_id: string;
    answer: string;
    pages: {
        document_id: string;
        title: string;
        page_number: number;
        page_image: string;
    }[];
    documents: RawDocument[];
    widgets: RawWidget[];
    is_partial: boolean;
};

type RawConversation = {
    id: string;
    title: string;
    created_at: string;
    updated_at: string;
};

type RawMessage = {
    id: string;
    role: "user" | "assistant";
    content: string;
    widgets: RawWidget[];
    documents: RawDocument[];
    is_partial: boolean;
};

// type GeneratedDocument = {
//     title: string;
//     theme: string;
//     outputPath: string;
// };

type ChatResult = {
    conversationId: string;
    title: string;
    answer: string;
    widgets: Widget[];
    documents: GeneratedDocument[];
    isPartial: boolean;
};

export type Conversation = {
    id: string;
    title: string;
    createdAt: string;
    updatedAt: string;
};

export type ChatMessage = {
    id: string;
    role: "user" | "assistant";
    content: string;
    widgets: Widget[];
    documents: GeneratedDocument[];
    isPartial: boolean;
};

export class ApiError extends Error {
    status: number;
    constructor(status: number, message: string) {
        super(message);
        this.status = status;
    }
}

async function ensureOk(res: Response): Promise<void> {
    if (res.ok) return;

    let message = `${res.status} ${res.statusText}`;
    try {
        const body = await res.json();
        if (body?.error) message = body.error;
    } catch {

    } throw new ApiError(res.status, message);
}

function mapWidgets(raw: RawWidget[] | null | undefined): Widget[] {
    return (raw ?? [])
        .map(parseWidget)
        .filter((w): w is Widget => w !== null);
}

function mapDocuments(raw: RawDocument[] | null | undefined): GeneratedDocument[] {
    return (raw ?? []).map((d) => ({
        title: d.title,
        theme: d.theme,
        outputPath: d.output_path,
        format: d.format ?? "",
    }));
}

export async function sendChatMessage(query: string, conversationId?: string | null, files: File[] = []): Promise<ChatResult> {
    const formData = new FormData();
    formData.append("query", query);
    if (conversationId) formData.append("conversation_id", conversationId);
    files.forEach((f) => formData.append("files", f));
    
    const res = await fetch(`${API_BASE_URL}/chat/agent`, {
        method: "POST",
        body: formData,
    });

    await ensureOk(res);

    const raw: ChatAgentResponse = await res.json();

    return {
        conversationId: raw.conversation_id,
        title: raw.title,
        answer: raw.answer,
        widgets: mapWidgets(raw.widgets),
        documents: mapDocuments(raw.documents),
        isPartial: raw.is_partial,
        skippedFiles: raw.skipped_files ?? [],
    };
}

export async function listConversations(): Promise<Conversation[]> {
    const res = await fetch (`${API_BASE_URL}/conversations`);
    await ensureOk(res);

    const raw: RawConversation[] = await res.json();

    return raw.map((c) => ({
        id: c.id,
        title: c.title,
        createdAt: c.created_at,
        updatedAt: c.updated_at,
    }));
}

export async function getConversationMessage(id: string): Promise<ChatMessage[]> {
    const res = await fetch (`${API_BASE_URL}/conversations/${id}/messages`);
    await ensureOk(res);

    const raw: RawMessage[] = await res.json();

    return raw.map((c) => ({
        id: c.id,
        role: c.role,
        content: c.content,
        widgets: mapWidgets(c.widgets),
        documents: mapDocuments(c.documents),
        isPartial: c.is_partial,
    }));
}

export async function renameConversation(id: string, title: string): Promise<Conversation> {
    const res = await fetch (`${API_BASE_URL}/conversations/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title }),
    });

    await ensureOk(res);

    const c: RawConversation = await res.json();
    return {
        id: c.id,
        title: c.title,
        createdAt: c.created_at,
        updatedAt: c.updated_at,
    };
}

export async function deleteConversation(id: string): Promise<void> {
    const res = await fetch (`${API_BASE_URL}/conversations/${id}`, {
        method: "DELETE",
    });
    await ensureOk(res);
}