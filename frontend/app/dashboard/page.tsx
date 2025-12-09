"use client";

import useSWR from 'swr';
import api from '@/lib/axios';
import { Package, Warehouse, AlertTriangle } from 'lucide-react';

const fetcher = (url: string) => api.get(url).then(res => res.data);

export default function DashboardPage() {
    const { data: items, error: itemsError } = useSWR('/items', fetcher);
    const { data: warehouses, error: whError } = useSWR('/warehouses', fetcher);

    const loading = !items && !itemsError;

    return (
        <div className="space-y-6">
            <h1 className="text-3xl font-bold text-gray-900">Dashboard Overview</h1>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Total Items Card */}
                <div className="bg-white p-6 rounded-xl shadow border border-gray-100 flex items-center justify-between">
                    <div>
                        <p className="text-gray-500 text-sm font-medium">Total Items</p>
                        <h3 className="text-3xl font-bold text-gray-900">{items?.length || 0}</h3>
                    </div>
                    <div className="p-3 bg-indigo-50 rounded-full text-indigo-600">
                        <Package className="w-8 h-8" />
                    </div>
                </div>

                {/* Warehouses Card */}
                <div className="bg-white p-6 rounded-xl shadow border border-gray-100 flex items-center justify-between">
                    <div>
                        <p className="text-gray-500 text-sm font-medium">Warehouses</p>
                        <h3 className="text-3xl font-bold text-gray-900">{warehouses?.length || 0}</h3>
                    </div>
                    <div className="p-3 bg-green-50 rounded-full text-green-600">
                        <Warehouse className="w-8 h-8" />
                    </div>
                </div>

                {/* Low Stock Alert */}
                <div className="bg-white p-6 rounded-xl shadow border border-gray-100 flex items-center justify-between">
                    <div>
                        <p className="text-gray-500 text-sm font-medium">Low Stock Alerts</p>
                        <h3 className="text-3xl font-bold text-gray-900">0</h3>
                    </div>
                    <div className="p-3 bg-red-50 rounded-full text-red-600">
                        <AlertTriangle className="w-8 h-8" />
                    </div>
                </div>
            </div>
        </div>
    );
}
