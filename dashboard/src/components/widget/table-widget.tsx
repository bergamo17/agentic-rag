import type { TableData } from "@/lib/widget-schema";

type TableWidgetProps = {
    title: string;
    data: TableData;
};

export function TableWidget({ title, data }: TableWidgetProps) {
    return (
        <div className="rounded-[4px] border border-gray-200 bg-white p-4">
            <p className="mb-3 text-sm text-gray-500">{title}</p>
            <table className="w-full border-collapse text-left">
                <thead>
                    <tr className="border-b border-gray-200">
                        {data.columns.map((col) => (
                            <th key={col} className="py-2 px-3 text-xs font-medium text-gray-500">
                                {col}
                            </th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {data.rows.map((row, i) => (
                        <tr
                            key={i}
                            className={`border-b border-gray-100 ${i % 2 === 1 ? "bg-gray-50/60" : ""}`}
                        >
                            {row.map((cell, j) => (
                                <td key={j} className="py-2 px-3 text-sm text-gray-900">
                                    {cell}
                                </td>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </table>
            <p className="mt-3 text-[10px] text-gray-400">Dibuat oleh AI · verifikasi sebelum digunakan</p>
        </div>
    );
}