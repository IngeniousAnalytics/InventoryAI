"use client";

import { useState } from 'react';
import useSWR from 'swr';
import api from '@/lib/axios';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

const fetcher = (url: string) => api.get(url).then(res => res.data);

export default function ItemsPage() {
    const { data: items, error, isLoading, mutate } = useSWR('/items', fetcher);

    const [saving, setSaving] = useState(false);
    const [actionError, setActionError] = useState('');

    const [newName, setNewName] = useState('');
    const [newSku, setNewSku] = useState('');
    const [newQty, setNewQty] = useState('0');
    const [newPrice, setNewPrice] = useState('0');

    const [editingId, setEditingId] = useState<string | null>(null);
    const [editName, setEditName] = useState('');
    const [editSku, setEditSku] = useState('');
    const [editQty, setEditQty] = useState('0');
    const [editPrice, setEditPrice] = useState('0');

    const createItem = async () => {
        setActionError('');
        if (!newName.trim()) {
            setActionError('Item name is required.');
            return;
        }
        const qty = Number(newQty);
        const price = Number(newPrice);
        if (Number.isNaN(qty) || Number.isNaN(price)) {
            setActionError('Quantity and price must be valid numbers.');
            return;
        }

        setSaving(true);
        try {
            await api.post('/items', {
                name: newName.trim(),
                sku: newSku.trim(),
                quantity: qty,
                price: price,
            });
            setNewName('');
            setNewSku('');
            setNewQty('0');
            setNewPrice('0');
            await mutate();
        } catch (e: any) {
            setActionError(e?.response?.data?.error || 'Failed to create item');
        } finally {
            setSaving(false);
        }
    };

    const startEdit = (item: any) => {
        setEditingId(item.id);
        setEditName(item.name || '');
        setEditSku(item.sku || '');
        setEditQty(String(item.quantity ?? 0));
        setEditPrice(String(item.price ?? 0));
    };

    const cancelEdit = () => {
        setEditingId(null);
        setEditName('');
        setEditSku('');
        setEditQty('0');
        setEditPrice('0');
    };

    const saveEdit = async () => {
        if (!editingId) return;
        setActionError('');
        if (!editName.trim()) {
            setActionError('Item name is required.');
            return;
        }
        const qty = Number(editQty);
        const price = Number(editPrice);
        if (Number.isNaN(qty) || Number.isNaN(price)) {
            setActionError('Quantity and price must be valid numbers.');
            return;
        }

        setSaving(true);
        try {
            await api.put(`/items/${editingId}`, {
                name: editName.trim(),
                sku: editSku.trim(),
                quantity: qty,
                price: price,
            });
            await mutate();
            cancelEdit();
        } catch (e: any) {
            setActionError(e?.response?.data?.error || 'Failed to update item');
        } finally {
            setSaving(false);
        }
    };

    const deleteItem = async (id: string) => {
        const ok = confirm('Delete this item?');
        if (!ok) return;
        setActionError('');
        setSaving(true);
        try {
            await api.delete(`/items/${id}`);
            await mutate();
            if (editingId === id) cancelEdit();
        } catch (e: any) {
            setActionError(e?.response?.data?.error || 'Failed to delete item');
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <h1 className="text-3xl font-bold text-gray-900">Inventory Items</h1>
            </div>

            <div className="bg-white rounded-lg shadow p-4">
                <div className="grid grid-cols-1 md:grid-cols-5 gap-3">
                    <Input placeholder="Name" value={newName} onChange={(e) => setNewName(e.target.value)} />
                    <Input placeholder="SKU" value={newSku} onChange={(e) => setNewSku(e.target.value)} />
                    <Input placeholder="Quantity" inputMode="numeric" value={newQty} onChange={(e) => setNewQty(e.target.value)} />
                    <Input placeholder="Price" inputMode="decimal" value={newPrice} onChange={(e) => setNewPrice(e.target.value)} />
                    <Button isLoading={saving} onClick={createItem}>Add Item</Button>
                </div>
                {actionError && <p className="mt-3 text-sm text-red-600">{actionError}</p>}
            </div>

            <div className="bg-white rounded-lg shadow overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Name
                                </th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    SKU
                                </th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Quantity
                                </th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Price
                                </th>
                                <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Actions
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {isLoading && (
                                <tr>
                                    <td colSpan={5} className="px-6 py-4 text-center text-sm text-gray-500">Loading items...</td>
                                </tr>
                            )}
                            {error && (
                                <tr>
                                    <td colSpan={5} className="px-6 py-4 text-center text-sm text-red-600">Failed to load items</td>
                                </tr>
                            )}
                            {items && items.length === 0 && (
                                <tr>
                                    <td colSpan={5} className="px-6 py-4 text-center text-sm text-gray-500">No items found. Add one to get started.</td>
                                </tr>
                            )}
                            {items?.map((item: any) => (
                                <tr key={item.id} className="hover:bg-gray-50 transition">
                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                        {editingId === item.id ? (
                                            <Input value={editName} onChange={(e) => setEditName(e.target.value)} />
                                        ) : (
                                            item.name
                                        )}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                        {editingId === item.id ? (
                                            <Input value={editSku} onChange={(e) => setEditSku(e.target.value)} />
                                        ) : (
                                            item.sku
                                        )}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                        {editingId === item.id ? (
                                            <Input inputMode="numeric" value={editQty} onChange={(e) => setEditQty(e.target.value)} />
                                        ) : (
                                            item.quantity
                                        )}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                        {editingId === item.id ? (
                                            <Input inputMode="decimal" value={editPrice} onChange={(e) => setEditPrice(e.target.value)} />
                                        ) : (
                                            `$${item.price}`
                                        )}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-right">
                                        <div className="flex justify-end gap-2">
                                            {editingId === item.id ? (
                                                <>
                                                    <Button isLoading={saving} onClick={saveEdit}>Save</Button>
                                                    <Button className="bg-gray-200 text-black hover:bg-gray-300" onClick={cancelEdit}>Cancel</Button>
                                                </>
                                            ) : (
                                                <>
                                                    <Button onClick={() => startEdit(item)}>Edit</Button>
                                                    <Button className="bg-red-600 hover:bg-red-700" isLoading={saving} onClick={() => deleteItem(item.id)}>
                                                        Delete
                                                    </Button>
                                                </>
                                            )}
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}
