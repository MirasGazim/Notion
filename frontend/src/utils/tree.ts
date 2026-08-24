export const buildBlockTree = (blocks: any[], parentId: string | null = null): any[] => {
  return blocks
    .filter((block) => block.parent_id === parentId)
    .sort((a, b) => a.position - b.position)
    .map((block) => ({
      ...block,
      children: buildBlockTree(blocks, block.id),
    }));
};
