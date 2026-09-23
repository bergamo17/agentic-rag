"use client";

import { useWidgetStore } from "@/lib/store/widgets-store";
import { WidgetGrid } from "@/components/widget/widget-grid";

export default function DashboardPage() {
    const widgets = useWidgetStore((state) => state.widgets);

    return (
        <div className="p-4">
            {widgets.length > 0 && (
                <div className="mb-4 rounded-[4px] border border-amber-200 bg-amber-50 px-4 py-2 text-xs text-amber-700">
                    Data pada widget di bawah dihasilkan oleh AI dan dapat mengandung kesalahan. Periksa kembali sebelum digunakan untuk keputusan penting.
                </div>
            )}
            <WidgetGrid widgets={widgets} />
        </div>
    );
}