"use client";

import useSWR from 'swr';
import api from '@/lib/axios';
import { Plus, MapPin } from 'lucide-react';
import { Button } from '@/components/ui/button';

const fetcher = (url: string) => api.get(url).then(res => res.data);

export default function WarehousesPage() {
    const { data: warehouses, error, isLoading } = useSWR('/warehouses', fetcher);

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <h1 className="text-3xl font-bold text-gray-900">Warehouses</h1>
                <Button onClick={() => alert("Add Warehouse Placeholder")}>
                    <Plus className="mr-2 h-4 w-4" /> Add Warehouse
                </Button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {isLoading && <p>Loading warehouses...</p>}

                {warehouses?.map((wh: any) => (
                    <div key={wh.id} className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 flex flex-col gap-4">
                        <div className="flex items-start justify-between">
                            <div>
                                <h3 className="text-lg font-semibold text-gray-900">{wh.name}</h3>
                                <div className="flex items-center text-sm text-gray-500 mt-1">
                                    <MapPin className="w-4 h-4 mr-1" />
                                    {wh.location || "No address provided"}
                                </div>
                            </div>
                        </div>
                        <div className="pt-4 border-t border-gray-50 flex items-center justify-between text-sm">
                            <span className="text-gray-500">Created: {new Date(wh.created_at).toLocaleDateString()}</span>
                        </div>
                    </div>
                ))}
                {warehouses && warehouses.length === 0 && (
                    <div className="col-span-full text-center py-10 bg-gray-50 rounded-lg border border-dashed border-gray-300">
                        <p className="text-gray-500">No warehouses found.</p>
                    </div>
                )}
            </div>
        </div>
    );
}
