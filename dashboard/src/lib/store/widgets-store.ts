import { create } from "zustand";
import type { Widget } from "@/lib/widget-schema";

type WidgetState = {
    widgets: Widget[];
    addWidgets: (newWidgets: Widget[]) => void;
    clearWidgets: () => void;
};

export const useWidgetStore = create<WidgetState>((set) => ({
    widgets: [],

    addWidgets: (newWidgets) =>
        set((state) => ({
            widgets: [...state.widgets, ...newWidgets],
        })),

    clearWidgets: () => set({widgets: []}),
}));