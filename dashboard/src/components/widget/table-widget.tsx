import type { TableData } from "@/lib/widget-schema";

type TableWidgetProps = {
    title: string;
    data: TableData;
};

export function TableWidget({ title, data }: TableWidgetProps) {
    return (
        <div className="rounded-soft-lg border border-border bg-card p-4">
            <p className="mb-3 text-[17px] font-semibold text-foreground">{title}</p>
            <div className="overflow-x-auto">
                <table className="w-full border-collapse text-left tabular-nums">
                    <thead>
                        <tr className="bg-grey-100">
                            {data.columns.map((col) => (
                                <th key={col} className="py-2 px-3 text-xs font-medium text-grey-600">
                                    {col}
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {data.rows.map((row, i) => (
                            <tr key={i} className="border-b border-border hover:bg-light-active">
                                {row.map((cell, j) => (
                                    <td key={j} className="py-2 px-3 text-sm text-foreground">
                                        {cell}
                                    </td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
            <p className="mt-3 text-[10px] text-grey-400">Dibuat oleh AI · verifikasi sebelum digunakan</p>
        </div>
    );
}