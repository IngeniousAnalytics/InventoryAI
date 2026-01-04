"use client";

import { useState } from 'react';
import useSWR from 'swr';
import api from '@/lib/axios';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { MapPin } from 'lucide-react';

const fetcher = (url: string) => api.get(url).then(res => res.data);

export default function WarehousesPage() {
    const { data: warehouses, error, isLoading, mutate } = useSWR('/warehouses', fetcher);

    const [name, setName] = useState('');
    const [location, setLocation] = useState('');
    const [saving, setSaving] = useState(false);
    const [actionError, setActionError] = useState('');

    const [editingId, setEditingId] = useState<string | null>(null);
    const [editName, setEditName] = useState('');
    const [editLocation, setEditLocation] = useState('');

    const createWarehouse = async () => {
        setActionError('');
        if (!name.trim()) {
            setActionError('Warehouse name is required.');
            return;
        }
        setSaving(true);
        try {
            await api.post('/warehouses', { name: name.trim(), location: location.trim() });
            setName('');
            setLocation('');
            await mutate();
        } catch (e: any) {
            setActionError(e?.response?.data?.error || 'Failed to create warehouse');
        } finally {
            setSaving(false);
        }
    };

    const startEdit = (wh: any) => {
        setEditingId(wh.id);
        setEditName(wh.name || '');
        setEditLocation(wh.location || '');
    };

    const cancelEdit = () => {
        setEditingId(null);
        setEditName('');
        setEditLocation('');
    };

    const saveEdit = async () => {
        if (!editingId) return;
        setActionError('');
        if (!editName.trim()) {
            setActionError('Warehouse name is required.');
            return;
        }
        setSaving(true);
        try {
            await api.put(`/warehouses/${editingId}`, { name: editName.trim(), location: editLocation });
            await mutate();
            cancelEdit();
        } catch (e: any) {
            setActionError(e?.response?.data?.error || 'Failed to update warehouse');
        } finally {
            setSaving(false);
        }
    };

    const deleteWarehouse = async (id: string) => {
        const ok = confirm('Delete this warehouse?');
        if (!ok) return;
        setActionError('');
        setSaving(true);
        try {
            await api.delete(`/warehouses/${id}`);
            await mutate();
            if (editingId === id) cancelEdit();
        } catch (e: any) {
            setActionError(e?.response?.data?.error || 'Failed to delete warehouse');
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <h1 className="text-3xl font-bold text-gray-900">Warehouses</h1>
            </div>

            <div className="bg-white rounded-lg shadow p-4">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                    <Input placeholder="Warehouse name" value={name} onChange={(e) => setName(e.target.value)} />
                    <Input placeholder="Location (optional)" value={location} onChange={(e) => setLocation(e.target.value)} />
                    <Button isLoading={saving} onClick={createWarehouse}>Add Warehouse</Button>
                </div>
                {actionError && <p className="mt-3 text-sm text-red-600">{actionError}</p>}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {isLoading && <p>Loading warehouses...</p>}
                {error && <p className="text-red-600">Failed to load warehouses</p>}

                {warehouses?.map((wh: any) => (
                    <div key={wh.id} className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 flex flex-col gap-4">
                        <div className="flex items-start justify-between">
                            <div>
                                {editingId === wh.id ? (
                                    <div className="space-y-2">
                                        <Input value={editName} onChange={(e) => setEditName(e.target.value)} />
                                        <Input value={editLocation} onChange={(e) => setEditLocation(e.target.value)} placeholder="Location (optional)" />
                                    </div>
                                ) : (
                                    <>
                                        <h3 className="text-lg font-semibold text-gray-900">{wh.name}</h3>
                                        <div className="flex items-center text-sm text-gray-500 mt-1">
                                            <MapPin className="w-4 h-4 mr-1" />
                                            {wh.location || "No address provided"}
                                        </div>
                                    </>
                                )}
                            </div>
                        </div>
                        <div className="pt-4 border-t border-gray-50 flex items-center justify-between text-sm">
                            <span className="text-gray-500">Created: {new Date(wh.created_at).toLocaleDateString()}</span>
                            <div className="flex gap-2">
                                {editingId === wh.id ? (
                                    <>
                                        <Button isLoading={saving} onClick={saveEdit}>Save</Button>
                                        <Button className="bg-gray-200 text-black hover:bg-gray-300" onClick={cancelEdit}>Cancel</Button>
                                    </>
                                ) : (
                                    <>
                                        <Button onClick={() => startEdit(wh)}>Edit</Button>
                                        <Button className="bg-red-600 hover:bg-red-700" isLoading={saving} onClick={() => deleteWarehouse(wh.id)}>
                                            Delete
                                        </Button>
                                    </>
                                )}
                            </div>
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
