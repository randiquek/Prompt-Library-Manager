import { useState, useEffect } from "react";
import { getAllPrompts } from "../services/api";
import './Prompts.css';

export default function Prompts({ onLogout }) {
    const [prompts, setPrompts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    const username = localStorage.getItem('username');

    useEffect(() => {
        fetchPrompts();
    }, []);

    const fetchPrompts = async () => {
        try {
            setLoading(true);
            const data = await getAllPrompts();
            setPrompts(data || []);
        } catch (err) {
            setError('Failed to load prompts');
            console.error(err)
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="prompts-container">
      <header className="prompts-header">
        <div>
          <h1>📝 Prompt Library</h1>
          <p>Welcome back, {username}!</p>
        </div>
        <button onClick={onLogout} className="logout-btn">
          Logout
        </button>
      </header>

      <div className="prompts-content">
        <div className="prompts-actions">
          <button className="add-btn">+ Add Prompt</button>
        </div>

        {loading && <div className="loading">Loading prompts...</div>}
        
        {error && <div className="error-message">{error}</div>}

        {!loading && !error && (
          <div className="prompts-table-container">
            <table className="prompts-table">
              <thead>
                <tr>
                  <th>Title</th>
                  <th>Category</th>
                  <th>Content Preview</th>
                  <th>Created</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {prompts.length === 0 ? (
                  <tr>
                    <td colSpan="5" className="no-data">
                      No prompts yet. Click "Add Prompt" to create one!
                    </td>
                  </tr>
                ) : (
                  prompts.map((prompt) => (
                    <tr key={prompt.id}>
                      <td className="prompt-title">{prompt.title}</td>
                      <td>
                        <span className="category-badge">{prompt.category || 'Uncategorized'}</span>
                      </td>
                      <td className="prompt-content">
                        {prompt.content.substring(0, 80)}...
                      </td>
                      <td className="prompt-date">
                        {new Date(prompt.created_at).toLocaleDateString()}
                      </td>
                      <td className="prompt-actions">
                        <button className="btn-edit">Edit</button>
                        <button className="btn-delete">Delete</button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
    );

}