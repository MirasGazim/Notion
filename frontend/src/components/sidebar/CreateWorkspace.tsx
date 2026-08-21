import { useState } from 'react';
import { authService } from '../../services/authService';
import '../../styles/layout/main.css';

export const CreateWorkspace = ({ onCreated }: { onCreated: () => void }) => {
  const [name, setName] = useState('');
  const [isOpen, setIsOpen] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await authService.createWorkspace({ name });
      setIsOpen(false);
      setName('');
      onCreated();
    } catch (error) {
      console.error('Ошибка создания:', error);
      alert('Ошибка при создании рабочего пространства');
    }
  };

  if (!isOpen) {
    return (
      <button 
        onClick={() => setIsOpen(true)} 
        className="sidebar-item" 
        style={{ marginTop: '16px', color: '#767671', border: 'none', background: 'none', width: '100%', textAlign: 'left' }}
      >
        + Add a workspace
      </button>
    );
  }

  return (
    <form onSubmit={handleSubmit} style={{ padding: '8px' }}>
      <input 
        value={name}
        onChange={(e) => setName(e.target.value)}
        placeholder="Workspace name"
        required
        style={{ width: '100%', padding: '4px', marginBottom: '4px' }}
      />
      <div style={{ display: 'flex', gap: '4px' }}>
        <button type="submit" style={{ flex: 1, padding: '4px' }}>Create</button>
        <button type="button" onClick={() => setIsOpen(false)} style={{ flex: 1, padding: '4px' }}>Cancel</button>
      </div>
    </form>
  );
};
