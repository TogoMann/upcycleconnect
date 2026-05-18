import { defineStore } from 'pinia';
import { useAuthStore } from './auth';
import { API_BASE } from '@/config';

export interface Plan {
    id: number;
    name: string;
    description: string;
    price: number;
    billing_cycle: string;
    features: string[];
    is_active: boolean;
}

export const usePlansStore = defineStore('plans', () => {
    const authStore = useAuthStore();

    async function getPlans() {
        const res = await fetch(`${API_BASE}/plans`);
        if (!res.ok) throw new Error('Failed to fetch plans');
        return await res.json();
    }

    async function choosePlan(planId: number, siret?: string) {
        const res = await fetch(`${API_BASE}/subscriptions/choose`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authStore.token}`
            },
            body: JSON.stringify({
                plan_id: planId,
                siret: siret
            })
        });
        if (!res.ok) {
            const error = await res.text();
            throw new Error(error || 'Failed to update plan');
        }
        await authStore.fetchCurrentUser(); // Refresh user role
    }

    return { getPlans, choosePlan };
});
