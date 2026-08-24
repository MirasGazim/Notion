import { useState } from 'react';
import { authService } from '../../services/authService';
import '../../styles/layout/main.css';

interface BlockProps {
  block: any;
  onUpdate: () => void;
}

export const Block = ({ block, onUpdate }: BlockProps) => {
  const [content, setContent] = useState(block.content.text || '');

  const handleBlur = async () => {
    // Автосохранение при потере фокуса
    await authService.updateBlock(block.id, { 
        type: block.type, 
        content: { text: content } 
    });
    onUpdate();
  };

  return (
    <div style={{ marginLeft: block.parent_id ? '20px' : '0', marginTop: '8px' }}>
      <input 
        value={content}
        onChange={(e) => setContent(e.target.value)}
        onBlur={handleBlur}
        style={{
            width: '100%',
            padding: '4px',
            border: 'none',
            fontSize: block.type === 'header' ? '24px' : '16px',
            outline: 'none'
        }}
      />
      
      {/* РЕКУРСИЯ: Рендерим детей */}
      {block.children && block.children.map((child: any) => (
        <Block key={child.id} block={child} onUpdate={onUpdate} />
      ))}
    </div>
  );
};
