// Iris Go 框架学习项目 - JavaScript 工具函数

// 工具函数命名空间
const IrisUtils = {
    // 初始化函数
    init() {
        this.bindEvents();
        this.initTooltips();
        this.initModals();
    },

    // 绑定全局事件
    bindEvents() {
        // 为所有外部链接添加安全属性
        document.addEventListener('DOMContentLoaded', () => {
            const externalLinks = document.querySelectorAll('a[href^="http"]:not([target])');
            externalLinks.forEach(link => {
                link.setAttribute('target', '_blank');
                link.setAttribute('rel', 'noopener noreferrer');
            });
        });

        // 表单提交确认
        document.addEventListener('submit', (e) => {
            const form = e.target;
            if (form.hasAttribute('data-confirm')) {
                const message = form.getAttribute('data-confirm');
                if (!confirm(message)) {
                    e.preventDefault();
                }
            }
        });

        // 删除确认
        document.addEventListener('click', (e) => {
            const element = e.target.closest('[data-delete-confirm]');
            if (element) {
                const message = element.getAttribute('data-delete-confirm');
                if (!confirm(message)) {
                    e.preventDefault();
                }
            }
        });
    },

    // 初始化工具提示
    initTooltips() {
        const tooltipElements = document.querySelectorAll('[data-tooltip]');
        tooltipElements.forEach(element => {
            element.addEventListener('mouseenter', (e) => {
                this.showTooltip(e.target);
            });
            
            element.addEventListener('mouseleave', (e) => {
                this.hideTooltip(e.target);
            });
        });
    },

    // 显示工具提示
    showTooltip(element) {
        const text = element.getAttribute('data-tooltip');
        if (!text) return;

        const tooltip = document.createElement('div');
        tooltip.className = 'iris-tooltip';
        tooltip.textContent = text;
        tooltip.style.cssText = `
            position: absolute;
            background: #333;
            color: white;
            padding: 5px 10px;
            border-radius: 4px;
            font-size: 12px;
            z-index: 1000;
            white-space: nowrap;
            pointer-events: none;
        `;

        document.body.appendChild(tooltip);

        const rect = element.getBoundingClientRect();
        tooltip.style.left = rect.left + (rect.width / 2) - (tooltip.offsetWidth / 2) + 'px';
        tooltip.style.top = rect.top - tooltip.offsetHeight - 5 + 'px';

        element._tooltip = tooltip;
    },

    // 隐藏工具提示
    hideTooltip(element) {
        if (element._tooltip) {
            element._tooltip.remove();
            delete element._tooltip;
        }
    },

    // 初始化模态框
    initModals() {
        const modalTriggers = document.querySelectorAll('[data-modal-target]');
        modalTriggers.forEach(trigger => {
            trigger.addEventListener('click', (e) => {
                e.preventDefault();
                const targetId = trigger.getAttribute('data-modal-target');
                this.showModal(targetId);
            });
        });

        // 关闭模态框
        document.addEventListener('click', (e) => {
            if (e.target.matches('[data-modal-close]') || e.target.matches('.modal-backdrop')) {
                this.hideModal();
            }
        });

        // ESC 键关闭模态框
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                this.hideModal();
            }
        });
    },

    // 显示模态框
    showModal(modalId) {
        const modal = document.getElementById(modalId);
        if (!modal) return;

        // 创建背景遮罩
        const backdrop = document.createElement('div');
        backdrop.className = 'modal-backdrop';
        backdrop.style.cssText = `
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0, 0, 0, 0.5);
            z-index: 999;
        `;

        modal.style.display = 'block';
        modal.style.position = 'fixed';
        modal.style.top = '50%';
        modal.style.left = '50%';
        modal.style.transform = 'translate(-50%, -50%)';
        modal.style.zIndex = '1000';
        modal.style.background = 'white';
        modal.style.padding = '20px';
        modal.style.borderRadius = '8px';
        modal.style.maxWidth = '90%';
        modal.style.maxHeight = '90%';
        modal.style.overflow = 'auto';

        document.body.appendChild(backdrop);
        this._currentModal = { modal, backdrop };
    },

    // 隐藏模态框
    hideModal() {
        if (this._currentModal) {
            this._currentModal.modal.style.display = 'none';
            this._currentModal.backdrop.remove();
            this._currentModal = null;
        }
    },

    // AJAX 请求封装
    ajax(options) {
        const defaults = {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            timeout: 30000
        };

        const config = Object.assign({}, defaults, options);

        return new Promise((resolve, reject) => {
            const xhr = new XMLHttpRequest();
            xhr.timeout = config.timeout;

            xhr.open(config.method, config.url, true);

            // 设置请求头
            Object.keys(config.headers).forEach(key => {
                xhr.setRequestHeader(key, config.headers[key]);
            });

            xhr.onload = () => {
                if (xhr.status >= 200 && xhr.status < 300) {
                    try {
                        const response = JSON.parse(xhr.responseText);
                        resolve(response);
                    } catch (e) {
                        resolve(xhr.responseText);
                    }
                } else {
                    reject(new Error(`HTTP ${xhr.status}: ${xhr.statusText}`));
                }
            };

            xhr.onerror = () => reject(new Error('网络错误'));
            xhr.ontimeout = () => reject(new Error('请求超时'));

            // 发送请求
            if (config.data) {
                if (typeof config.data === 'object') {
                    xhr.send(JSON.stringify(config.data));
                } else {
                    xhr.send(config.data);
                }
            } else {
                xhr.send();
            }
        });
    },

    // GET 请求
    get(url, options = {}) {
        return this.ajax(Object.assign({}, options, { method: 'GET', url }));
    },

    // POST 请求
    post(url, data, options = {}) {
        return this.ajax(Object.assign({}, options, { method: 'POST', url, data }));
    },

    // PUT 请求
    put(url, data, options = {}) {
        return this.ajax(Object.assign({}, options, { method: 'PUT', url, data }));
    },

    // DELETE 请求
    delete(url, options = {}) {
        return this.ajax(Object.assign({}, options, { method: 'DELETE', url }));
    },

    // 显示通知
    showNotification(message, type = 'info', duration = 3000) {
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 15px 20px;
            border-radius: 4px;
            color: white;
            font-weight: 500;
            z-index: 10000;
            max-width: 300px;
            word-wrap: break-word;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            transform: translateX(100%);
            transition: transform 0.3s ease;
        `;

        // 设置背景色
        const colors = {
            success: '#28a745',
            error: '#dc3545',
            warning: '#ffc107',
            info: '#17a2b8'
        };
        notification.style.backgroundColor = colors[type] || colors.info;

        notification.textContent = message;
        document.body.appendChild(notification);

        // 显示动画
        setTimeout(() => {
            notification.style.transform = 'translateX(0)';
        }, 100);

        // 自动隐藏
        setTimeout(() => {
            notification.style.transform = 'translateX(100%)';
            setTimeout(() => {
                if (notification.parentNode) {
                    notification.parentNode.removeChild(notification);
                }
            }, 300);
        }, duration);
    },

    // 格式化日期
    formatDate(date, format = 'YYYY-MM-DD HH:mm:ss') {
        if (!(date instanceof Date)) {
            date = new Date(date);
        }

        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        const seconds = String(date.getSeconds()).padStart(2, '0');

        return format
            .replace('YYYY', year)
            .replace('MM', month)
            .replace('DD', day)
            .replace('HH', hours)
            .replace('mm', minutes)
            .replace('ss', seconds);
    },

    // 相对时间
    timeAgo(date) {
        if (!(date instanceof Date)) {
            date = new Date(date);
        }

        const seconds = Math.floor((new Date() - date) / 1000);
        const intervals = {
            年: 31536000,
            月: 2592000,
            天: 86400,
            小时: 3600,
            分钟: 60
        };

        if (seconds < 5) return '刚刚';

        for (const [unit, secondsInUnit] of Object.entries(intervals)) {
            const interval = Math.floor(seconds / secondsInUnit);
            if (interval >= 1) {
                return `${interval}${unit}前`;
            }
        }

        return '刚刚';
    },

    // 防抖函数
    debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    },

    // 节流函数
    throttle(func, limit) {
        let inThrottle;
        return function(...args) {
            if (!inThrottle) {
                func.apply(this, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    },

    // 复制到剪贴板
    copyToClipboard(text) {
        if (navigator.clipboard) {
            return navigator.clipboard.writeText(text).then(() => {
                this.showNotification('已复制到剪贴板', 'success');
            });
        } else {
            // 降级方案
            const textArea = document.createElement('textarea');
            textArea.value = text;
            textArea.style.position = 'fixed';
            textArea.style.opacity = '0';
            document.body.appendChild(textArea);
            textArea.select();
            
            try {
                document.execCommand('copy');
                this.showNotification('已复制到剪贴板', 'success');
            } catch (err) {
                this.showNotification('复制失败', 'error');
            }
            
            document.body.removeChild(textArea);
        }
    },

    // 本地存储
    storage: {
        set(key, value) {
            try {
                localStorage.setItem(key, JSON.stringify(value));
            } catch (e) {
                console.error('设置本地存储失败:', e);
            }
        },

        get(key, defaultValue = null) {
            try {
                const item = localStorage.getItem(key);
                return item ? JSON.parse(item) : defaultValue;
            } catch (e) {
                console.error('获取本地存储失败:', e);
                return defaultValue;
            }
        },

        remove(key) {
            try {
                localStorage.removeItem(key);
            } catch (e) {
                console.error('删除本地存储失败:', e);
            }
        },

        clear() {
            try {
                localStorage.clear();
            } catch (e) {
                console.error('清空本地存储失败:', e);
            }
        }
    }
};

// API 客户端
const ApiClient = {
    baseURL: '/api',
    token: null,

    // 设置认证令牌
    setToken(token) {
        this.token = token;
        if (token) {
            IrisUtils.storage.set('auth_token', token);
        } else {
            IrisUtils.storage.remove('auth_token');
        }
    },

    // 获取认证令牌
    getToken() {
        if (!this.token) {
            this.token = IrisUtils.storage.get('auth_token');
        }
        return this.token;
    },

    // 发送请求
    async request(endpoint, options = {}) {
        const url = endpoint.startsWith('http') ? endpoint : this.baseURL + endpoint;
        const token = this.getToken();

        const config = {
            headers: {},
            ...options
        };

        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }

        try {
            const response = await IrisUtils.ajax({ ...config, url });
            return response;
        } catch (error) {
            console.error('API 请求失败:', error);
            throw error;
        }
    },

    // GET 请求
    get(endpoint, options = {}) {
        return this.request(endpoint, { ...options, method: 'GET' });
    },

    // POST 请求
    post(endpoint, data, options = {}) {
        return this.request(endpoint, { ...options, method: 'POST', data });
    },

    // PUT 请求
    put(endpoint, data, options = {}) {
        return this.request(endpoint, { ...options, method: 'PUT', data });
    },

    // DELETE 请求
    delete(endpoint, options = {}) {
        return this.request(endpoint, { ...options, method: 'DELETE' });
    }
};

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    IrisUtils.init();
    
    // 全局错误处理
    window.addEventListener('error', (e) => {
        console.error('全局错误:', e.error);
        IrisUtils.showNotification('页面出现错误，请刷新重试', 'error');
    });

    // 全局未处理的 Promise 拒绝
    window.addEventListener('unhandledrejection', (e) => {
        console.error('未处理的 Promise 拒绝:', e.reason);
        IrisUtils.showNotification('请求处理失败，请重试', 'error');
    });
});

// 导出到全局
window.IrisUtils = IrisUtils;
window.ApiClient = ApiClient;