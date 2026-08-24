import { useState, useMemo } from 'react';
import { MainLayout } from '../components/layout/MainLayout';
import { authService } from '../services/authService';
import { buildBlockTree } from '../utils/tree';
import { Block } from '../components/editor/Block';

export const MainPage = () => {
  const [selectedWorkspaceId, setSelectedWorkspaceId] = useState<string | null>(null);
  const [blocks, setBlocks] = useState([]);

  const handleWorkspaceSelect = async (id: string) => {
    setSelectedWorkspaceId(id);
    try {
      const response = await authService.getWorkspaceBlocks(id);
      setBlocks(response.data.blocks || []);
    } catch (error) {
      console.error('Ошибка загрузки блоков:', error);
      setBlocks([]);
    }
  };

  const blockTree = useMemo(() => buildBlockTree(blocks), [blocks]);

  return (
    <MainLayout onWorkspaceSelect={handleWorkspaceSelect}>
      {selectedWorkspaceId ? (
        <div>
          {blockTree.length === 0 ? (
            <div style={{ color: '#767671', marginTop: '20px' }}>
              <p>Здесь пока пусто. Создайте первый блок, чтобы начать работу.</p>
            </div>
          ) : (
            <div>
              {blockTree.map(block => (
                <Block key={block.id} block={block} onUpdate={() => handleWorkspaceSelect(selectedWorkspaceId!)} />
              ))}
            </div>
          )}
        </div>
      ) : (
        <h1>Выберите рабочее пространство из боковой панели</h1>
      )}
    </MainLayout>
  );
};
