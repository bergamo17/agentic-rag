import { Widget } from "@/lib/widget-schema";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { PanelRightOpen } from "lucide-react";
import { useWidgetStore } from "@/lib/store/widgets-store";

type GeneratedDocument = {
    title: string;
    theme: string;
    outputPath: string;
    format: "docx"|"pdf"|"xlsx"|"md";
}

type Message = {
    id: string;
    role: "user" | "agent";
    content: string;
    isPartial?: boolean;
    widgets?: Widget[];
    documents?: GeneratedDocument[];
};

type MessageBubbleProps = {
    message: Message;
}

export function MessageBubble({ message }: MessageBubbleProps){
    const isUser = message.role === "user";
    const openPanel = useWidgetStore((state) => state.openPanel);
    const widgets = message.widgets ?? [];

    return (
        <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
            <div
                className={`max-w-[75%] rounded-[10px] px-4 py-3 ${
                isUser 
                    ? "bg-blue-600 text-white" 
                    : "bg-white border border-gray-200 border-l-4 border-l-indigo-600 text-gray-900"
                }`}
            >
                {isUser ? (
                    <p className="text-sm whitespace-pre-wrap">{message.content}</p>
                ) : (
                    <div
                        className="prose prose-sm max-w-none
                                   prose-headings:mt-4 prose-headings:mb-2
                                   prose-h1:text-lg prose-h2:text-base prose-h3:text-sm
                                   prose-p:my-2 prose-li:my-0.5 prose-hr:my-4"
                    >
                        <ReactMarkdown
                            remarkPlugins={[remarkGfm]}
                            components={{
                                // tabel lebar bisa di-scroll horizontal, tidak melebarkan bubble
                                table: ({ children }) => (
                                    <div className="overflow-x-auto">
                                        <table>{children}</table>
                                    </div>
                                ),
                            }}
                        >
                            {message.content}
                        </ReactMarkdown>
                    </div>
                )}

                {widgets.length > 0 && (
                    <button
                        onClick={() => openPanel(widgets[0].id)}
                        className="mt-2 flex items-center gap-1.5 rounded-lg border border-gray-200 bg-gray-50 px-3 py-1.5 text-xs font-medium text-indigo-600 hover:bg-gray-100"
                    >
                        <PanelRightOpen size={14} />
                        Lihat Visualisasi ({widgets.length})
                    </button>
                )}

                {message.documents?.map((doc, i) => (
                    <div
                        key={i}
                        className="mt-2 flex items-center justify-between rounded-lg border border-gray-200 bg-gray-50 px-3 py-2"
                    >
                        <div className="min-w-0">
                            <p className="truncate text-sm font-medium text-gray-900">{doc.title}</p>
                            <p className="text-xs text-gray-500">{doc.theme}</p>
                        </div>
                            <a
                            href={`${process.env.NEXT_PUBLIC_API_BASE_URL}/documents/download?path=${encodeURIComponent(doc.outputPath)}`}
                            download
                            className="ml-3 shrink-0 text-xs font-medium text-indigo-600 underline"
                        >
                            Download
                        </a>
                    </div>
                ))}
                
                {message.isPartial && (
                <span className="mt-1 inline-block text-xs text-amber-600">
                    ⚠ Jawaban belum lengkap
                </span>
                )}
            </div>
        </div>
    );
}