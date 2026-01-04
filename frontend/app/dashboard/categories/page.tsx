"use client";

import { useMemo, useState } from "react";
import useSWR from "swr";
import api from "@/lib/axios";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type Category = {
    id: string;
    name: string;
    created_at?: string;
};

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export default function CategoriesPage() {
    const { data: categories, error, isLoading, mutate } = useSWR<Category[]>("/categories", fetcher);

    const [newName, setNewName] = useState("");
    const [saving, setSaving] = useState(false);
    const [actionError, setActionError] = useState<string>("");

    const [editingId, setEditingId] = useState<string | null>(null);
    const [editName, setEditName] = useState("");

    const sorted = useMemo(() => {
        if (!categories) return [];
        return [...categories].sort((a, b) => (a.name || "").localeCompare(b.name || ""));
    }, [categories]);

    const createCategory = async () => {
        setActionError("");
        if (!newName.trim()) {
            setActionError("Category name is required.");
            return;
        }
        setSaving(true);
        try {
            await api.post("/categories", { name: newName.trim() });
            setNewName("");
            await mutate();
        } catch (e: any) {
            setActionError(e?.response?.data?.error || "Failed to create category");
        } finally {
            setSaving(false);
        }
    };

    const startEdit = (cat: Category) => {
        setEditingId(cat.id);
        setEditName(cat.name || "");
    };

    const cancelEdit = () => {
        setEditingId(null);
        setEditName("");
    };

    const saveEdit = async () => {
        if (!editingId) return;
        setActionError("");
        if (!editName.trim()) {
            setActionError("Category name is required.");
            return;
        }
        setSaving(true);
        try {
            await api.put(`/categories/${editingId}`, { name: editName.trim() });
            await mutate();
            cancelEdit();
        } catch (e: any) {
            setActionError(e?.response?.data?.error || "Failed to update category");
        } finally {
            setSaving(false);
        }
    };

    const deleteCategory = async (id: string) => {
        const ok = confirm("Delete this category?");
        if (!ok) return;
        setActionError("");
        setSaving(true);
        try {
            await api.delete(`/categories/${id}`);
            await mutate();
            if (editingId === id) cancelEdit();
        } catch (e: any) {
            setActionError(e?.response?.data?.error || "Failed to delete category");
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <h1 className="text-3xl font-bold text-gray-900">Categories</h1>
            </div>

            <div className="bg-white rounded-lg shadow p-4">
                <div className="flex flex-col md:flex-row gap-3 md:items-center">
                    <div className="flex-1">
                        <Input
                            placeholder="New category name"
                            value={newName}
                            onChange={(e) => setNewName(e.target.value)}
                        />
                    </div>
                    <Button isLoading={saving} onClick={createCategory}>
                        Add Category
                    </Button>
                </div>
                {actionError && <p className="mt-3 text-sm text-red-600">{actionError}</p>}
            </div>

            <div className="bg-white rounded-lg shadow overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th
                                    scope="col"
                                    className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                >
                                    Name
                                </th>
                                <th
                                    scope="col"
                                    className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"
                                >
                                    Actions
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {isLoading && (
                                <tr>
                                    <td colSpan={2} className="px-6 py-4 text-center text-sm text-gray-500">
                                        Loading categories...
                                    </td>
                                </tr>
                            )}
                            {error && (
                                <tr>
                                    <td colSpan={2} className="px-6 py-4 text-center text-sm text-red-600">
                                        Failed to load categories
                                    </td>
                                </tr>
                            )}
                            {sorted.length === 0 && !isLoading && !error && (
                                <tr>
                                    <td colSpan={2} className="px-6 py-4 text-center text-sm text-gray-500">
                                        No categories yet.
                                    </td>
                                </tr>
                            )}

                            {sorted.map((cat) => {
                                const isEditing = editingId === cat.id;
                                return (
                                    <tr key={cat.id} className="hover:bg-gray-50 transition">
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {isEditing ? (
                                                <Input value={editName} onChange={(e) => setEditName(e.target.value)} />
                                            ) : (
                                                cat.name
                                            )}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right">
                                            <div className="flex justify-end gap-2">
                                                {isEditing ? (
                                                    <>
                                                        <Button isLoading={saving} onClick={saveEdit}>
                                                            Save
                                                        </Button>
                                                        <Button className="bg-gray-200 text-black hover:bg-gray-300" onClick={cancelEdit}>
                                                            Cancel
                                                        </Button>
                                                    </>
                                                ) : (
                                                    <>
                                                        <Button onClick={() => startEdit(cat)}>Edit</Button>
                                                        <Button
                                                            className="bg-red-600 hover:bg-red-700"
                                                            isLoading={saving}
                                                            onClick={() => deleteCategory(cat.id)}
                                                        >
                                                            Delete
                                                        </Button>
                                                    </>
                                                )}
                                            </div>
                                        </td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}
