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

function isStringArray(v: unknown): v is string[] {
    return Array.isArray(v) && v.every((item) => typeof item === "string");
}

function isNumberArray(v: unknown): v is number[] {
    return Array.isArray(v) && v.every((item) => typeof item === "number");
}

export function parseWidget(raw: {
    widget_type: string;
    title: string;
    data: string;
}): Widget | null {
    try {
        const parsedData = JSON.parse(raw.data);

        switch (raw.widget_type) {
            case "chart":
                if (
                    !Array.isArray(parsedData.labels) ||
                    !Array.isArray(parsedData.values) ||
                    parsedData.values.length !== parsedData.labels.length ||
                    !parsedData.values.every((v: unknown) => typeof v === "number") ||
                    (parsedData.chartType !== "bar" && parsedData.chartType !== "line")
                ) {
                    return null;
                }
                return { widget_type: "chart", title: raw.title, data: parsedData as ChartData };

            case "table":
                if (
                    !isStringArray(parsedData.columns) ||
                    !Array.isArray(parsedData.rows) ||
                    !parsedData.rows.every(
                        (row: unknown) => isStringArray(row) && row.length === parsedData.columns.length
                    )
                ) {
                    return null;
                }
                return { widget_type: "table", title: raw.title, data: parsedData as TableData };

            case "card":
                if (
                    typeof parsedData.value !== "string" ||
                    (parsedData.description !== undefined && typeof parsedData.description !== "string")
                ) {
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