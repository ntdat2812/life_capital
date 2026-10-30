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

export const KNOWN_GOLD_BRANDS = [
  'DOJI',
  'Bảo Tín Minh Châu',
  'Bảo Tín Mạnh Hải',
  'SJC',
  'Phú Quý',
  'PNJ',
  'Mi Hồng',
  'Tư nhân',
  'Khác'
]

export const getGoldBrand = (asset) => {
  if (!asset) return 'Khác'
  const name = asset.name || ''
  const ticker = asset.ticker || ''
  
  if (name.includes(' - ')) {
    const prefix = name.split(' - ')[0].trim()
    for (const b of KNOWN_GOLD_BRANDS) {
      if (prefix.toLowerCase() === b.toLowerCase()) return b
    }
  }
  
  if (ticker.startsWith('GOLD_')) {
    const parts = ticker.split('_')
    if (parts.length >= 2) {
      const code = parts[1]
      if (code === 'BTMC') return 'Bảo Tín Minh Châu'
      if (code === 'BTMH') return 'Bảo Tín Mạnh Hải'
      if (code === 'DOJI') return 'DOJI'
      if (code === 'SJC') return 'SJC'
      if (code === 'PNJ') return 'PNJ'
      if (code === 'PHUQUY') return 'Phú Quý'
      if (code === 'MIHONG') return 'Mi Hồng'
    }
  }
  
  for (const b of KNOWN_GOLD_BRANDS) {
    if (name.toLowerCase().includes(b.toLowerCase())) return b
  }
  
  return 'Khác'
}

export const getBrandBadgeClass = (_brand) => {
  return 'bg-amber-500/15 text-amber-300 border-amber-500/30'
}

