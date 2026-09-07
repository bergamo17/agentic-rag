export type ChartData = {
    chartType: "bar" | "line";
    labels: string[];
    values: number[];
};

export type TableData = {
    columns: string[];
    rows: string[][];
};

export type CardData = {
    value: string;
    description?: string;
};

export type Widget = 
    | { widget_type: "chart"; title: string; data: ChartData }
    | { widget_type: "table"; title: string; data:TableData }
    | { widget_type: "card"; title: string; data:CardData }

export function parseWidget(raw: {
    widget_type: string;
    title: string;
    data: string;
}): Widget | null {
    try {
        const parsedData = JSON.parse(raw.data);

        switch (raw.widget_type) {
            case "chart":
                if (!Array.isArray(parsedData.labels) || !Array.isArray(parsedData.values)) {
                    return null;
                }
                return { widget_type: "chart", title: raw.title, data: parsedData as ChartData };

            case "table":
                if (!Array.isArray(parsedData.columns) || !Array.isArray(parsedData.rows)) {
                    return null;
                }
                return { widget_type: "table", title: raw.title, data: parsedData as TableData };

            case "card":
                if (typeof parsedData.value !== "string") {
                    return null;
                }
                return { widget_type: "card", title: raw.title, data: parsedData as CardData };

            default:
                return null;
        }
    } catch {
        return null;
    }
}