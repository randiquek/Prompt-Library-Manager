import { useState } from "react";
import Modal from "./Modal";

export default function PromptForm({ isOpen, onClose, onSubmit, initialData = null }) {
    const [title, setTitle] = useState(initialData?.title || '');
    const [content, setContent] = useState(initialData?.content || '');
    const [category, setCategory] = useState(initialData?.category || '');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        if (!title.trim()) {
            setError('Title is required');
            return;
        }
        if (!content.trim()) {
            setError('Content is required');
            return;
        }

        setLoading(true);

        try {
            await onSubmit({ title, content, category });
            // form reset
            setTitle('');
            setContent('');
            setCategory('');
            onClose();
        } catch (err) {
            setError(err.response?.data || 'Failed to save prompt');
        } finally {
            setLoading(false);
        }

    };

}