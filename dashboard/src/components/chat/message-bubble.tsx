import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { PanelRightOpen, Paperclip, TriangleAlert } from "lucide-react";
import { useWidgetStore } from "@/lib/store/widgets-store";
import type { Message } from "@/lib/chat-type";

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
                className={`max-w-[75%] rounded-soft-lg px-4 py-3 ${
                isUser
                    ? "bg-primary text-white"
                    : "bg-card border border-border border-l-4 border-l-primary text-foreground"
                }`}
            >
                {isUser ? (
                    <>
                        {message.content && (
                            <p className="text-sm whitespace-pre-wrap">{message.content}</p>
                        )}
                        {message.attachments && message.attachments.length > 0 && (
                            <div className="mt-2 flex flex-wrap gap-1.5">
                                {message.attachments.map((name, i) => (
                                    <span
                                        key={`${name}-${i}`}
                                        className="flex items-center gap-1 rounded-tight-sm bg-white/15 px-2 py-0.5 text-xs"
                                    >
                                        <Paperclip size={11} />
                                        {name}
                                    </span>
                                ))}
                            </div>
                        )}
                    </>
                ) : (
                    <div
                        className="prose prose-sm max-w-none
                                   prose-headings:mt-4 prose-headings:mb-2 prose-headings:text-foreground
                                   prose-h1:text-lg prose-h2:text-base prose-h3:text-sm
                                   prose-p:my-2 prose-p:text-foreground prose-li:my-0.5 prose-li:text-foreground
                                   prose-strong:text-foreground prose-a:text-primary-accent prose-hr:my-4"
                    >
                        <ReactMarkdown
                            remarkPlugins={[remarkGfm]}
                            components={{
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
                        className="mt-2 flex items-center gap-1.5 rounded-md border border-border bg-grey-50 px-3 py-1.5 text-xs font-medium text-primary-accent hover:bg-grey-100"
                    >
                        <PanelRightOpen size={14} />
                        Lihat Visualisasi ({widgets.length})
                    </button>
                )}

                {message.documents?.map((doc, i) => (
                    <div
                        key={i}
                        className="mt-2 flex items-center justify-between rounded-md border border-border bg-grey-50 px-3 py-2"
                    >
                        <div className="min-w-0">
                            <p className="truncate text-sm font-medium text-grey-900">{doc.title}</p>
                            <p className="text-xs text-grey-600">{doc.theme}</p>
                        </div>
                        <a
                            href={`${process.env.NEXT_PUBLIC_API_BASE_URL}/documents/download?path=${encodeURIComponent(doc.outputPath)}`}
                            download
                            className="ml-3 shrink-0 text-xs font-medium text-primary-accent underline"
                        >
                            Download
                        </a>
                    </div>
                ))}

                {message.isPartial && (
                    <span className="mt-2 inline-flex items-center gap-1 rounded-full bg-warning-soft px-2 py-0.5 text-xs font-medium text-warning-accent">
                        <TriangleAlert size={12} />
                        Jawaban belum lengkap
                    </span>
                )}
            </div>
        </div>
    );
}