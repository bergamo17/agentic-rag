import { create } from "zustand";
import type { Widget } from "@/lib/widget-schema";

type WidgetState = {
    widgets: Widget[];
    isPanelOpen: boolean;
    activeWidgetId: string | null;

    addWidgets: (newWidgets: Widget[]) => void;
    clearWidgets: () => void;
    closePanel: () => void;
    openPanel: (widgetId?: string) => void;
};

export const useWidgetStore = create<WidgetState>((set) => ({
    widgets: [],
    isPanelOpen: false,
    activeWidgetId: null,

    addWidgets: (newWidgets) =>
        set((state) => ({
            widgets: [...state.widgets, ...newWidgets],
            isPanelOpen: true,
            activeWidgetId: newWidgets.at(-1)?.id ?? state.activeWidgetId,
        })),

    clearWidgets: () => set({ widgets: [], isPanelOpen: false, activeWidgetId: null }),

    closePanel: () => set({ isPanelOpen: false }),

    openPanel: (widgetId) =>
        set((state) => ({
            isPanelOpen: true,
            activeWidgetId: widgetId ?? state.activeWidgetId,
        })),
}));