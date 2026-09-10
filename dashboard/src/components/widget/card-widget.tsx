import type { CardData } from "@/lib/widget-schema";

type CardWidgetProps = {
    title: string;
    data: CardData;
};

export function CardWidget({ title, data }: CardWidgetProps ) {
    return (
        <div className="rounded-[4px] border border-gray-200 bg-white p-4">
            <p className="text-sm text-gray-500">{title}</p>
            <p className="mt-2 font-mono text-3xl font-semibold text-gray-900">
                {data.value}
            </p>
            {data.description && (
                <p className="mt-1 text-xs text-gray-500">{data.description}</p>
            )}
            <p className="mt-3 text-[10px] text-gray-400">Dibuat oleh AI · verifikasi sebelum digunakan</p>
        </div>
    );
}