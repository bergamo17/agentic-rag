type Message = {
    id: string;
    role: "user" | "agent";
    content: string;
    isPartial?: boolean;
};

type MessageBubbleProps = {
    message: Message;
}

export function MessageBubble({ message }: MessageBubbleProps){
    const isUser = message.role === "user";

    return (
        <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
            <div
                className={`max-w-[75%] rounded-[10px] px-4 py-3 ${
                isUser 
                    ? "bg-blue-600 text-white" 
                    : "bg-white border border-gray-200 border-l-4 border-l-indigo-600 text-gray-900"
                }`}
            >
                <p className="text-sm whitespace-pre-wrap">{message.content}</p>
                {message.isPartial && (
                <span className="mt-1 inline-block text-xs text-amber-600">
                    ⚠ Jawaban belum lengkap
                </span>
                )}
            </div>
        </div>
    );
}