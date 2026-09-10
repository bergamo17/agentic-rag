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

export function ChartWidget({ title, data }: ChartWidgetProps) {
    const points = data.labels.map((label, i) => ({
        label,
        value: data.values[i],
    }));

    return (
        <div className="rounded-[4px] border border-gray-200 bg-white p-4">
            <p className="mb-3 text-sm text-gray-500">{title}</p>
            <ResponsiveContainer width="100%" height={220}>
                {data.chartType === "bar" ? (
                    <BarChart data={points}>
                        <CartesianGrid vertical={false} stroke="#E4E7EC" />
                        <XAxis dataKey="label" tick={{ fontSize: 12, fill: "#6B7280" }} axisLine={false} tickLine={false} />
                        <YAxis tick={{ fontSize: 12, fill: "#6B7280" }} axisLine={false} tickLine={false} />
                        <Tooltip contentStyle={{ borderRadius: 4, borderColor: "#E4E7EC" }} />
                        <Bar dataKey="value" fill="#0EA5A4" radius={[2, 2, 0, 0]} />
                    </BarChart>
                ) : (
                    <LineChart data={points}>
                        <CartesianGrid vertical={false} stroke="#E4E7EC" />
                        <XAxis dataKey="label" tick={{ fontSize: 12, fill: "#6B7280" }} axisLine={false} tickLine={false} />
                        <YAxis tick={{ fontSize: 12, fill: "#6B7280" }} axisLine={false} tickLine={false} />
                        <Tooltip contentStyle={{ borderRadius: 4, borderColor: "#E4E7EC" }} />
                        <Line type="monotone" dataKey="value" stroke="#0EA5A4" strokeWidth={2} dot={false} />
                    </LineChart>
                )}
            </ResponsiveContainer>
            <p className="mt-3 text-[10px] text-gray-400">Dibuat oleh AI · verifikasi sebelum digunakan</p>
        </div>
    );
}