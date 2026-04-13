import toast from 'react-hot-toast';

const toastStyles = {
  style: {
    background: 'linear-gradient(135deg, #f6d365 0%, #fda085 100%)',
    color: 'white',
    fontWeight: 500,
    padding: '16px',
    borderRadius: '24px',
  },
  iconTheme: {
    primary: 'white',
    secondary: '#E65100',
  },
};

export const showToast = {
  success: (message: string) => toast.success(message, toastStyles),
  error: (message: string) => toast.error(message, toastStyles),
  loading: (message: string) => toast.loading(message, toastStyles),
  ai: (message: string) => toast(message, {
    ...toastStyles,
    icon: '🤖',
  }),
};