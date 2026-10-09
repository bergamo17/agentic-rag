import type { CardData } from "@/lib/widget-schema";

type CardWidgetProps = {
    title: string;
    data: CardData;
};

export function CardWidget({ title, data }: CardWidgetProps ) {
    return (
        <div className="rounded-soft-lg border border-border bg-card p-4">
            <p className="text-sm text-muted">{title}</p>
            <p className="mt-2 font-mono text-3xl font-semibold text-foreground">
                {data.value}
            </p>
            {data.description && (
                <p className="mt-1 text-xs text-muted">{data.description}</p>
            )}
            <p className="mt-3 text-[10px] text-grey-400">Dibuat oleh AI · verifikasi sebelum digunakan</p>
        </div>
    );
}