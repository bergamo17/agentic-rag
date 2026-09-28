"use client";

import { useState } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';

const SIDEBAR_WIDTH = 260;
const SIDEBAR_COLLAPSED_WIDTH = 72;

export default function Sidebar() {
    const [collapsed, setCollapsed] = useState(false);

    return (
        <div
            className="relative h-screen border-r bg-white transition-all duration-300 ease-in-out"
            style={{ width: collapsed ? SIDEBAR_COLLAPSED_WIDTH : SIDEBAR_WIDTH }}
        >
            {/* Toggle button — nempel di tepi kanan sidebar, vertically centered */}
            <button
                onClick={() => setCollapsed((prev) => !prev)}
                className="absolute -right-3 top-1/2 -translate-y-1/2 z-10
                           flex h-6 w-6 items-center justify-center
                           rounded-full border bg-white shadow-md
                           hover:bg-gray-50 transition-colors"
                aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
            >
                {collapsed ? (
                    <ChevronRight className="h-4 w-4 text-gray-600" />
                ) : (
                    <ChevronLeft className="h-4 w-4 text-gray-600" />
                )}
            </button>

            {/* Isi sidebar */}
            <div className="flex flex-col h-full p-4 overflow-hidden">
                <h1
                    className={`font-bold text-lg mb-6 whitespace-nowrap transition-opacity duration-200 ${
                        collapsed ? "opacity-0" : "opacity-100"
                    }`}
                >
                    Agentic RAG
                </h1>

                <nav className="flex flex-col gap-1">
                    {/* contoh item nav — ganti sesuai punya kamu */}
                    <SidebarItem icon="chat" label="Chat" collapsed={collapsed} href="/main/chat" />
                    <SidebarItem icon="dashboard" label="Dashboard" collapsed={collapsed} href="/main/dashboard" />
                </nav>
            </div>
        </div>
    );
}

function SidebarItem({ icon, label, collapsed, href }: {
    icon: React.ReactNode;
    label: string;
    collapsed: boolean;
    href: string;
}) {
    return (
        <a
            href={href}
            className="flex items-center gap-3 rounded-lg px-3 py-2 hover:bg-gray-100 transition-colors"
        >
            <span className="shrink-0">{icon}</span>
            {!collapsed && <span className="whitespace-nowrap">{label}</span>}
        </a>
    );
}