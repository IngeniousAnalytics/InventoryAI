"use client";

import Link from "next/link";
import { useAuthStore } from "@/store/useAuthStore";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { LayoutDashboard, Package, Warehouse, LogOut } from "lucide-react";

export default function DashboardLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    const { token, logout } = useAuthStore();
    const router = useRouter();

    useEffect(() => {
        if (!token) {
            router.push("/auth/login");
        }
    }, [token, router]);

    if (!token) return null;

    return (
        <div className="flex h-screen bg-gray-100">
            {/* Sidebar */}
            <aside className="w-64 bg-white shadow-md">
                <div className="p-6 border-b">
                    <h1 className="text-2xl font-bold text-indigo-600">InventoryAI</h1>
                </div>
                <nav className="p-4 space-y-2">
                    <Link
                        href="/dashboard"
                        className="flex items-center gap-3 px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 rounded-lg transition"
                    >
                        <LayoutDashboard className="w-5 h-5" />
                        Dashboard
                    </Link>
                    <Link
                        href="/dashboard/items"
                        className="flex items-center gap-3 px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 rounded-lg transition"
                    >
                        <Package className="w-5 h-5" />
                        Items
                    </Link>
                    <Link
                        href="/dashboard/warehouses"
                        className="flex items-center gap-3 px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 rounded-lg transition"
                    >
                        <Warehouse className="w-5 h-5" />
                        Warehouses
                    </Link>
                </nav>
                <div className="p-4 border-t mt-auto">
                    <button
                        onClick={() => {
                            logout();
                            router.push("/auth/login");
                        }}
                        className="flex w-full items-center gap-3 px-4 py-3 text-red-600 hover:bg-red-50 rounded-lg transition"
                    >
                        <LogOut className="w-5 h-5" />
                        Sign Out
                    </button>
                </div>
            </aside>

            {/* Main Content */}
            <main className="flex-1 overflow-y-auto p-8">
                {children}
            </main>
        </div>
    );
}
