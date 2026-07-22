export const ASSET_CATEGORIES = {
  cash: { name: 'Tiền mặt', icon: '💵' },
  deposit: { name: 'Tiền gửi', icon: '🏦' },
  gold: { name: 'Vàng', icon: '✨' }, // Gold bar / Sparkles
  stock: { name: 'Cổ phiếu', icon: '📈' },
  fund: { name: 'Chứng chỉ quỹ', icon: '📦' },
  crypto: { name: 'Tiền điện tử', icon: '₿' },
  real_estate: { name: 'Bất động sản', icon: '🏢' }
}

export const LIABILITY_CATEGORIES = {
  mortgage: { name: 'Vay mua nhà', icon: '🏘️' },
  auto_loan: { name: 'Vay mua xe', icon: '🚗' },
  student_loan: { name: 'Vay học tập', icon: '🎓' },
  credit_card: { name: 'Thẻ tín dụng', icon: '💳' },
  personal_loan: { name: 'Vay tiêu dùng', icon: '👤' },
  other: { name: 'Khác', icon: '🧾' }
}

export const getCategoryName = (key) => ASSET_CATEGORIES[key]?.name || key
export const getCategoryIcon = (key) => ASSET_CATEGORIES[key]?.icon || '💰'

export const getLiabilityName = (key) => LIABILITY_CATEGORIES[key]?.name || key
export const getLiabilityIcon = (key) => LIABILITY_CATEGORIES[key]?.icon || '💸'
