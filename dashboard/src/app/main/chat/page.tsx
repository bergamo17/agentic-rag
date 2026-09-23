"use client";

import { ChatWindow } from "@/components/chat/chat-window";
import { DashboardPanel } from "@/components/dashboard/dashboard-panel";
import { useWidgetStore } from "@/lib/store/widgets-store";

export default function ChatPage() {
    const isPanelOpen = useWidgetStore((state) => state.isPanelOpen);

    return (
        <div className="flex h-full overflow-hidden">
            <div
                className={`h-full transition-all duration-300 ease-in-out ${
                    isPanelOpen ? "w-[38%]" : "w-full"
                }`}
            >
                <ChatWindow/>
            </div>
            {isPanelOpen && (
                <div className="h-full flex-1">
                    <DashboardPanel />
                </div>
            )}
        </div>
    );
}