import { parseWidget, Widget } from "@/lib/widget-schema";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

type ChatAgentResponse = {
    answer: string;
    pages: {
        document_id: string;
        title: string;
        page_number: number;
        page_image: string;
    }[];
    documents: {
        title: string;
        theme: string;
        output_path: string;
    }[];
    widgets: {
        widget_type: string;
        title: string;
        data: string;
    }[];
    is_partial: boolean;
};

type GeneratedDocument = {
    title: string;
    theme: string;
    outputPath: string;
}

type ChatResult = {
    answer: string;
    widgets: Widget[];
    documents: GeneratedDocument[];
    isPartial: boolean;
}

export async function sendChatMessage(query: string): Promise<ChatResult> {
    const res = await fetch(`${API_BASE_URL}/chat/agent`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query }),
    });

    if (!res.ok) {
        throw new Error(`API error: ${res.status} ${res.statusText}`);
    }

    const raw: ChatAgentResponse = await res.json();

    const parsedWidget = raw.widgets
        .map(parseWidget)
        .filter((w): w is Widget => w !== null);

    const documents: GeneratedDocument[] = raw.documents.map((d) => ({
        title: d.title,
        theme: d.theme,
        outputPath: d.output_path,
    }));

    return {
        answer: raw.answer,
        widgets: parsedWidget,
        documents,
        isPartial: raw.is_partial,
    };
}