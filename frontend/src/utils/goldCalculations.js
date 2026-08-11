export const getGoldMultiplier = (unit) => {
  if (unit === 'Loại 1 Lượng' || unit === 'Lượng') return 1;
  if (unit === 'Loại 5 Chỉ') return 0.5;
  if (unit === 'Loại 2 Chỉ') return 0.2;
  if (unit === 'Loại 1 Chỉ' || unit === 'Chỉ') return 0.1;
  if (unit === 'Loại 0.5 Chỉ') return 0.05;
  if (unit === 'Phân') return 0.01;
  return 1;
};

export const extractGoldUnit = (name) => {
  const match = name?.match(/\((.+?)\)$/);
  return match ? match[1] : 'Lượng';
};
