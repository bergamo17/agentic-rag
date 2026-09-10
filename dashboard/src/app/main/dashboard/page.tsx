"use client";

import { useWidgetStore } from "@/lib/store/widgets-store";
import { CardWidget } from "@/components/widget/card-widget";
import { TableWidget } from "@/components/widget/table-widget";
import { ChartWidget } from "@/components/widget/chart-widget";

export default function DashboardPage() {
    const widgets = useWidgetStore((state) => state.widgets);

    if (widgets.length === 0) {
        return (
            <div className="flex h-full items-center justify-center p-4">
                <p className="text-sm text-gray-500">
                    Belum ada widget. Tanyakan sesuatu yang menghasilkan data di halaman Chat untuk melihatnya di sini.
                </p>
            </div>
        );
    }

    return (
        <div className="p-4">
            <div className="mb-4 rounded-[4px] border border-amber-200 bg-amber-50 px-4 py-2 text-xs text-amber-700">
                Data pada widget di bawah dihasilkan oleh AI dan dapat mengandung kesalahan. Periksa kembali sebelum digunakan untuk keputusan penting.
            </div>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
                {widgets.map((widget, i) => {
                    switch (widget.widget_type) {
                        case "card":
                            return <CardWidget key={i} title={widget.title} data={widget.data} />;
                        case "table":
                            return <TableWidget key={i} title={widget.title} data={widget.data} />;
                        case "chart":
                            return <ChartWidget key={i} title={widget.title} data={widget.data} />;
                        default:
                            return null;
                    }
                })}
            </div>
        </div>
    );
}