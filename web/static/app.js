const { createApp } = Vue;
const { ElMessage, ElMessageBox } = ElementPlus;

// 配置axios
const api = axios.create({
    baseURL: '/api/v1',
    timeout: 10000
});

// 请求拦截器 - 自动添加token
api.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    error => {
        return Promise.reject(error);
    }
);

// 响应拦截器 - 处理认证错误
api.interceptors.response.use(
    response => {
        return response;
    },
    error => {
        if (error.response?.status === 401) {
            // token过期或无效，清除本地存储并跳转到登录页
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/login.html';
        }
        return Promise.reject(error);
    }
);

const app = createApp({
    data() {
        return {
            // 系统监控数据
            systemStats: {},
            statsLoading: false,
            
            // 任务数据
            tasks: [],
            tasksLoading: false,
            taskPagination: {
                page: 1,
                pageSize: 10,
                total: 0
            },
            
            // 搜索和筛选
            searchKeyword: '',
            statusFilter: '',
            
            // 用户信息
            currentUser: null,
            token: null,
            
            // 任务表单
            showTaskDialog: false,
            editingTask: null,
            taskForm: {
                name: '',
                schedule: '',
                command: '',
                method: 'GET',
                headers: '',
                description: '',
                enabled: true
            },
            taskRules: {
                name: [
                    { required: true, message: '请输入任务名称', trigger: 'blur' }
                ],
                schedule: [
                    { required: true, message: '请输入调度表达式', trigger: 'blur' }
                ],
                command: [
                    { required: true, message: '请输入命令或URL', trigger: 'blur' }
                ]
            },
            saving: false,
            
            // 任务日志
            showLogsDialog: false,
            taskLogs: [],
            logsLoading: false,
            currentLogTaskId: null,
            logPagination: {
                page: 1,
                pageSize: 10,
                total: 0
            }
        };
    },
    
    computed: {
        showHttpConfig() {
            return this.taskForm.command.startsWith('http://') || 
                   this.taskForm.command.startsWith('https://');
        },
        
        // 过滤后的任务列表
        filteredTasks() {
            let filtered = this.tasks;
            
            // 关键词搜索
            if (this.searchKeyword) {
                const keyword = this.searchKeyword.toLowerCase();
                filtered = filtered.filter(task => 
                    task.name.toLowerCase().includes(keyword) ||
                    task.command.toLowerCase().includes(keyword) ||
                    (task.description && task.description.toLowerCase().includes(keyword))
                );
            }
            
            // 状态筛选
            if (this.statusFilter) {
                if (this.statusFilter === 'enabled') {
                    filtered = filtered.filter(task => task.enabled);
                } else if (this.statusFilter === 'disabled') {
                    filtered = filtered.filter(task => !task.enabled);
                }
            }
            
            return filtered;
        }
    },
    
    mounted() {
        // 检查登录状态
        this.checkAuth();
        
        this.loadSystemStats();
        this.loadTasks();
        
        // 定时刷新系统监控数据
        setInterval(() => {
            this.loadSystemStats();
        }, 30000); // 30秒刷新一次
    },
    
    methods: {
        // 加载系统监控数据
        async loadSystemStats() {
            this.statsLoading = true;
            try {
                const response = await api.get('/system/stats');
                this.systemStats = response.data;
            } catch (error) {
                console.error('加载系统监控数据失败:', error);
                ElMessage.error('加载系统监控数据失败');
            } finally {
                this.statsLoading = false;
            }
        },
        
        // 加载任务列表
        async loadTasks() {
            this.tasksLoading = true;
            try {
                const response = await api.get('/tasks/paginated', {
                    params: {
                        page: this.taskPagination.page,
                        page_size: this.taskPagination.pageSize
                    }
                });
                this.tasks = response.data.data || [];
                this.taskPagination.total = response.data.total || 0;
            } catch (error) {
                console.error('加载任务列表失败:', error);
                ElMessage.error('加载任务列表失败');
            } finally {
                this.tasksLoading = false;
            }
        },
        
        // 任务分页处理
        handleTaskPageChange(page) {
            this.taskPagination.page = page;
            this.loadTasks();
        },
        
        handleTaskPageSizeChange(pageSize) {
            this.taskPagination.pageSize = pageSize;
            this.taskPagination.page = 1;
            this.loadTasks();
        },
        
        // 编辑任务
        editTask(task) {
            this.editingTask = task;
            this.taskForm = {
                name: task.name,
                schedule: task.schedule,
                command: task.command,
                method: task.method || 'GET',
                headers: task.headers || '',
                description: task.description || '',
                enabled: task.enabled
            };
            this.showTaskDialog = true;
        },
        
        // 保存任务
        async saveTask() {
            try {
                await this.$refs.taskFormRef.validate();
                
                this.saving = true;
                const taskData = {
                    Name: this.taskForm.name,
                    Schedule: this.taskForm.schedule,
                    Command: this.taskForm.command,
                    Method: this.showHttpConfig ? this.taskForm.method : '',
                    Headers: this.showHttpConfig ? this.taskForm.headers : '',
                    Description: this.taskForm.description,
                    Enabled: this.taskForm.enabled
                };
                
                if (this.editingTask) {
                    await api.put(`/tasks/${this.editingTask.id}`, taskData);
                    ElMessage.success('任务更新成功');
                } else {
                    await api.post('/tasks', taskData);
                    ElMessage.success('任务创建成功');
                }
                
                this.showTaskDialog = false;
                this.resetTaskForm();
                this.loadTasks();
            } catch (error) {
                console.error('保存任务失败:', error);
                ElMessage.error('保存任务失败: ' + (error.response?.data?.error || error.message));
            } finally {
                this.saving = false;
            }
        },
        
        // 删除任务
        async deleteTask(id) {
            try {
                await ElMessageBox.confirm('确定要删除这个任务吗？', '确认删除', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                });
                
                await api.delete(`/tasks/${id}`);
                ElMessage.success('任务删除成功');
                this.loadTasks();
            } catch (error) {
                if (error !== 'cancel') {
                    console.error('删除任务失败:', error);
                    ElMessage.error('删除任务失败');
                }
            }
        },
        
        // 启用/禁用任务
        async toggleTask(task) {
            try {
                const newStatus = !task.enabled;
                await api.put(`/tasks/${task.id}`, {
                    Name: task.name,
                    Schedule: task.schedule,
                    Command: task.command,
                    Method: task.method,
                    Headers: task.headers,
                    Description: task.description,
                    Enabled: newStatus
                });
                
                ElMessage.success(`任务${newStatus ? '启用' : '禁用'}成功`);
                this.loadTasks();
            } catch (error) {
                console.error('切换任务状态失败:', error);
                ElMessage.error('操作失败');
            }
        },
        
        // 立即执行任务
        async executeTask(id) {
            try {
                await api.post(`/tasks/${id}/execute`);
                ElMessage.success('任务执行成功');
            } catch (error) {
                console.error('执行任务失败:', error);
                ElMessage.error('执行任务失败: ' + (error.response?.data?.error || error.message));
            }
        },
        
        // 显示任务日志
        async showTaskLogs(task) {
            this.currentLogTaskId = task.id;
            this.logPagination.page = 1;
            this.showLogsDialog = true;
            await this.loadTaskLogs();
        },
        
        // 加载任务日志
        async loadTaskLogs() {
            if (!this.currentLogTaskId) return;
            
            this.logsLoading = true;
            try {
                const response = await api.get(`/tasks/${this.currentLogTaskId}/logs/paginated`, {
                    params: {
                        page: this.logPagination.page,
                        page_size: this.logPagination.pageSize
                    }
                });
                this.taskLogs = response.data.data || [];
                this.logPagination.total = response.data.total || 0;
            } catch (error) {
                console.error('加载任务日志失败:', error);
                ElMessage.error('加载任务日志失败');
            } finally {
                this.logsLoading = false;
            }
        },
        
        // 日志分页处理
        handleLogPageChange(page) {
            this.logPagination.page = page;
            this.loadTaskLogs();
        },
        
        handleLogPageSizeChange(pageSize) {
            this.logPagination.pageSize = pageSize;
            this.logPagination.page = 1;
            this.loadTaskLogs();
        },
        
        // 重置任务表单
        resetTaskForm() {
            this.editingTask = null;
            this.taskForm = {
                name: '',
                schedule: '',
                command: '',
                method: 'GET',
                headers: '',
                description: '',
                enabled: true
            };
            if (this.$refs.taskFormRef) {
                this.$refs.taskFormRef.resetFields();
            }
        },
        
        // 命令输入变化处理
        onCommandChange() {
            // 当命令变为URL时，自动设置默认的HTTP方法
            if (this.showHttpConfig && !this.taskForm.method) {
                this.taskForm.method = 'GET';
            }
        },
        
        // 格式化时间
        formatTime(timestamp) {
            if (!timestamp) return '-';
            return new Date(timestamp).toLocaleString('zh-CN');
        },
        
        // 格式化字节数
        formatBytes(bytes) {
            if (!bytes || bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
        },
        
        // 格式化运行时间
        formatUptime(seconds) {
            if (!seconds) return '0秒';
            const days = Math.floor(seconds / 86400);
            const hours = Math.floor((seconds % 86400) / 3600);
            const minutes = Math.floor((seconds % 3600) / 60);
            
            if (days > 0) {
                return `${days}天 ${hours}小时`;
            } else if (hours > 0) {
                return `${hours}小时 ${minutes}分钟`;
            } else {
                return `${minutes}分钟`;
            }
        },
        
        // 获取CPU颜色
        getCPUColor(percentage) {
            if (percentage < 50) return '#10b981';
            if (percentage < 80) return '#f59e0b';
            return '#ef4444';
        },
        
        // 获取内存颜色
        getMemoryColor(percentage) {
            if (percentage < 60) return '#10b981';
            if (percentage < 85) return '#f59e0b';
            return '#ef4444';
        },
        
        // 获取磁盘颜色
        getDiskColor(percentage) {
            if (percentage < 70) return '#10b981';
            if (percentage < 90) return '#f59e0b';
            return '#ef4444';
        },
        
        // 检查认证状态
        checkAuth() {
            this.token = localStorage.getItem('token');
            const userStr = localStorage.getItem('user');
            
            if (!this.token) {
                window.location.href = '/login.html';
                return;
            }
            
            if (userStr) {
                this.currentUser = JSON.parse(userStr);
            }
        },
        

        
        // 登出
        logout() {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/login.html';
        },

        // 格式化字节数
        formatBytes(bytes) {
            if (!bytes || bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        },

        // 格式化运行时间
        formatUptime(seconds) {
            if (!seconds || seconds === 0) return '0秒';
            
            const days = Math.floor(seconds / 86400);
            const hours = Math.floor((seconds % 86400) / 3600);
            const minutes = Math.floor((seconds % 3600) / 60);
            const secs = seconds % 60;
            
            let result = '';
            if (days > 0) result += `${days}天 `;
            if (hours > 0) result += `${hours}小时 `;
            if (minutes > 0) result += `${minutes}分钟 `;
            if (secs > 0 || result === '') result += `${secs}秒`;
            
            return result.trim();
        },
        
        // 搜索处理
        handleSearch() {
            // 搜索时重置到第一页
            this.taskPagination.page = 1;
        },
        
        // 状态筛选处理
        handleStatusFilter() {
            // 筛选时重置到第一页
            this.taskPagination.page = 1;
        },
        
        // 重置筛选条件
        resetFilters() {
            this.searchKeyword = '';
            this.statusFilter = '';
            this.taskPagination.page = 1;
        }
    },
    
    watch: {
        showTaskDialog(val) {
            if (!val) {
                this.resetTaskForm();
            }
        }
    }
});

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component);
}

app.use(ElementPlus);
app.mount('#app');