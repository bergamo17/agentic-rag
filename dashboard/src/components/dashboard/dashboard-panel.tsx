"use client";

import { X } from "lucide-react";
import { useWidgetStore } from "@/lib/store/widgets-store";
import { WidgetGrid } from "../widget/widget-grid";

export function DashboardPanel() {
    const widgets = useWidgetStore((state) => state.widgets);
    const closePanel = useWidgetStore((state) => state.closePanel);

    return (
        <div className="flex h-full flex-col border-l border-gray-200 bg-white">
            <div className="flex items-center justify-between border-b border-gray-200 px-4 py-2">
                <span className="text-sm font-medium text-gray-900">Visualisasi</span>
                <button onClick={closePanel} aria-lambel="Tutup Panel" className="text-gray-400 hover:text-gray-600">
                    <X size={16} />
                </button>
            </div>
            <div className="flex-1 overflow-y-auto">
                <WidgetGrid widgets={widgets} columns="grid-cols-1" />
            </div>
        </div>
    );
}