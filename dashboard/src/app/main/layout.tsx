"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { MessageSquare, LayoutGrid } from "lucide-react";

const navItems = [
  {href: "/main/chat", label: "Chat", icon: MessageSquare},
  {href: "/main/dashboard", label: "Dashboard", icon: LayoutGrid},
];

export default function MainLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  return(
    <div className="flex h-screen">
      <aside className="w-56 shrink-0 border-r border-gray-200 flex flex-col gap-1 p-4 bg-white">
        <div className="mb-4 px-3 text-sm font-semibold text-gray-900">
          Agentic RAG
        </div>
        {navItems.map(({ href, label, icon: Icon }) => {
          const isActive = pathname === href;
          return (
            <Link
              key={href}
              href={href}
              className={`flex items-center gap-2 rounded-[10px] px-3 py-2 text-sm font-medium transition-colors ${
                isActive
                  ? "bg-indigo-50 text-indigo-600"
                  : "text-gray-600 hover:bg-gray-50"
              }`}
            >
              <Icon size={16} />
              {label}
            </Link>
          );
        })}
      </aside>
      <main className="flex-1 overflow-hidden bg-gray-50">{children}</main>
    </div>
  );
}