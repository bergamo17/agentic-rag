"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { MessageSquare, LayoutGrid, ChevronLeft, ChevronRight } from "lucide-react";

const navItems = [
  { href: "/main/chat", label: "Chat", icon: MessageSquare },
  { href: "/main/dashboard", label: "Dashboard", icon: LayoutGrid },
];

export default function MainLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="flex h-screen">
      <aside
        className={`relative shrink-0 border-r border-gray-200 flex flex-col gap-1 p-4 bg-white transition-all duration-300 ease-in-out ${
          collapsed ? "w-16" : "w-56"
        }`}
      >
        {/* Toggle button */}
        <button
          onClick={() => setCollapsed((prev) => !prev)}
          className="absolute -right-3 top-6 z-10 flex h-6 w-6 items-center justify-center
                     rounded-full border border-gray-200 bg-white shadow-sm
                     hover:bg-gray-50 transition-colors"
          aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
          {collapsed ? (
            <ChevronRight size={14} className="text-gray-600" />
          ) : (
            <ChevronLeft size={14} className="text-gray-600" />
          )}
        </button>

        <div
          className={`mb-4 px-3 text-sm font-semibold text-gray-900 whitespace-nowrap overflow-hidden transition-opacity duration-200 ${
            collapsed ? "opacity-0" : "opacity-100"
          }`}
        >
          Agentic RAG
        </div>

        {navItems.map(({ href, label, icon: Icon }) => {
          const isActive = pathname === href;
          return (
            <Link
              key={href}
              href={href}
              title={collapsed ? label : undefined}
              className={`flex items-center gap-2 rounded-[10px] px-3 py-2 text-sm font-medium transition-colors ${
                isActive
                  ? "bg-indigo-50 text-indigo-600"
                  : "text-gray-600 hover:bg-gray-50"
              }`}
            >
              <Icon size={16} className="shrink-0" />
              {!collapsed && <span className="whitespace-nowrap">{label}</span>}
            </Link>
          );
        })}
      </aside>
      <main className="flex-1 overflow-hidden bg-gray-50">{children}</main>
    </div>
  );
}