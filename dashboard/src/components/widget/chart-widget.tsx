"use client";

import {
    BarChart,
    Bar,
    LineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
} from "recharts";
import type { ChartData } from "@/lib/widget-schema";

type ChartWidgetProps = {
    title: string;
    data: ChartData;
};

const CHART_1 = "#0468CC";
const GRID = "#E3E9F0";
const TICK = { fontSize: 12, fill: "#647186", fontFamily: "var(--font-geist-mono)" };
const TOOLTIP_STYLE = { borderRadius: 4, borderColor: GRID, fontSize: 12 };

export function ChartWidget({ title, data }: ChartWidgetProps) {
    const points = data.labels.map((label, i) => ({
        label,
        value: data.values[i],
    }));

    return (
        <div className="rounded-soft-lg border border-border bg-card p-4">
            <p className="mb-3 text-[17px] font-semibold text-foreground">{title}</p>
            <ResponsiveContainer width="100%" height={220}>
                {data.chartType === "bar" ? (
                    <BarChart data={points}>
                        <CartesianGrid vertical={false} stroke={GRID} />
                        <XAxis dataKey="label" tick={TICK} axisLine={false} tickLine={false} />
                        <YAxis tick={TICK} axisLine={false} tickLine={false} />
                        <Tooltip contentStyle={TOOLTIP_STYLE} />
                        <Bar dataKey="value" fill={CHART_1} radius={[3, 3, 0, 0]} />
                    </BarChart>
                ) : (
                    <LineChart data={points}>
                        <CartesianGrid vertical={false} stroke={GRID} />
                        <XAxis dataKey="label" tick={TICK} axisLine={false} tickLine={false} />
                        <YAxis tick={TICK} axisLine={false} tickLine={false} />
                        <Tooltip contentStyle={TOOLTIP_STYLE} />
                        <Line type="monotone" dataKey="value" stroke={CHART_1} strokeWidth={2} dot={false} />
                    </LineChart>
                )}
            </ResponsiveContainer>
            <p className="mt-3 text-[10px] text-grey-400">Dibuat oleh AI · verifikasi sebelum digunakan</p>
        </div>
    );
}