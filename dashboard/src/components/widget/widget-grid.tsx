"use client";

import { ChartWidget } from "./chart-widget";
import { CardWidget } from "./card-widget";
import { TableWidget } from "./table-widget";
import type { Widget } from "@/lib/widget-schema";

export function WidgetGrid ({
    widgets,
    columns = "grid-cols-1 md:grid-cols-2 lg:grid-cols-3",
}: {
    widgets: Widget[],
    columns?: string,
}) {
    if (widgets.length === 0) {
        return (
            <div className="flex h-full items-center justify-center p-4">
                <p className="text-sm text-muted">
                    Belum ada widget. Tanyakan sesuatu yang menghasilkan data untuk melihatnya disini.
                </p>
            </div>
        );
    }

    return (
        <div className={`grid gap-4 p-4 ${columns}`}>
            {widgets.map((widget, i) => {
                switch (widget.widget_type) {
                    case "card":
                        return <CardWidget key={widget.id ?? i} title={widget.title} data={widget.data} />;
                    case "table":
                        return <TableWidget key={i} title={widget.title} data={widget.data} />;
                    case "chart":
                        return <ChartWidget key={i} title={widget.title} data={widget.data} />;
                    default:
                        return null;
                }
            })}
        </div>
    );
}