// 全局变量
const API_BASE = '/api/v1';
let editingTaskId = null;

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    // 绑定事件监听器
    document.getElementById('refreshTasks').addEventListener('click', loadTasks);
    document.getElementById('refreshStats').addEventListener('click', loadSystemStats);
    document.getElementById('taskForm').addEventListener('submit', handleTaskSubmit);
    document.getElementById('cancelEdit').addEventListener('click', cancelEdit);
    document.getElementById('taskCommand').addEventListener('input', function() {
        // 根据命令是否为URL显示HTTP配置
        toggleHttpConfig();
    });
    
    // 初始加载数据
    loadTasks();
    loadSystemStats();
    
    // 初始化时根据命令是否为URL显示HTTP配置
    toggleHttpConfig();
});

// 加载任务列表
function loadTasks() {
    fetch(`${API_BASE}/tasks`)
        .then(response => response.json())
        .then(tasks => {
            const tbody = document.querySelector('#taskTable tbody');
            tbody.innerHTML = '';
            
            tasks.forEach(task => {
                const row = tbody.insertRow();
                row.innerHTML = `
                    <td>${task.id}</td>
                    <td>${task.name}</td>
                    <td>${task.schedule}</td>
                    <td>${task.command}</td>
                    <td>${task.method || 'N/A'}</td>
                    <td>${task.enabled ? '启用' : '禁用'}</td>
                    <td>${task.description || ''}</td>
                    <td>
                        <button onclick="editTask(${task.id})">编辑</button>
                        <button onclick="deleteTask(${task.id})">删除</button>
                        <button onclick="toggleTask(${task.id}, ${!task.enabled})">${task.enabled ? '禁用' : '启用'}</button>
                        <button onclick="executeTask(${task.id})">立即执行</button>
                        <button onclick="showTaskLogs(${task.id})">查看日志</button>
                    </td>
                `;
            });
        })
        .catch(error => {
            console.error('加载任务失败:', error);
            alert('加载任务失败');
        });
}

// 加载系统监控数据
function loadSystemStats() {
    fetch(`${API_BASE}/system/stats`)
        .then(response => response.json())
        .then(stats => {
            const statsDiv = document.getElementById('systemStats');
            statsDiv.innerHTML = `
                <div class="stat-item">
                    <span class="stat-label">CPU使用率:</span>
                    <span class="stat-value">${stats.CPUUsage.toFixed(2)}%</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">内存使用率:</span>
                    <span class="stat-value">${stats.MemoryUsage.toFixed(2)}%</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">系统负载:</span>
                    <span class="stat-value">${stats.SystemLoad.toFixed(2)}</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">更新时间:</span>
                    <span class="stat-value">${new Date(stats.Timestamp).toLocaleString()}</span>
                </div>
            `;
        })
        .catch(error => {
            console.error('加载系统监控数据失败:', error);
            alert('加载系统监控数据失败');
        });
}

// 编辑任务
function editTask(id) {
    fetch(`${API_BASE}/tasks/${id}`)
        .then(response => response.json())
        .then(task => {
            document.getElementById('taskId').value = task.id;
            document.getElementById('taskName').value = task.name;
            document.getElementById('taskSchedule').value = task.schedule;
            document.getElementById('taskCommand').value = task.command;
            document.getElementById('taskMethod').value = task.method || 'GET';
            document.getElementById('taskHeaders').value = task.headers || '';
            document.getElementById('taskDescription').value = task.description || '';
            document.getElementById('taskEnabled').checked = task.enabled;
            
            editingTaskId = id;
            
            // 根据任务是否有HTTP方法设置HTTP配置显示选项
            if (task.method) {
                document.getElementById('showHttpConfig').value = 'show';
            } else {
                document.getElementById('showHttpConfig').value = 'auto';
            }
            
            // 根据命令是否为URL显示HTTP配置
            toggleHttpConfig();
        })
        .catch(error => {
            console.error('加载任务详情失败:', error);
            alert('加载任务详情失败');
        });
}

// 删除任务
function deleteTask(id) {
    if (!confirm('确定要删除这个任务吗？')) {
        return;
    }
    
    fetch(`${API_BASE}/tasks/${id}`, {
        method: 'DELETE'
    })
        .then(response => {
            if (response.ok) {
                loadTasks();
                alert('任务删除成功');
            } else {
                throw new Error('删除失败');
            }
        })
        .catch(error => {
            console.error('删除任务失败:', error);
            alert('删除任务失败');
        });
}

// 启用/禁用任务
function toggleTask(id, enable) {
    const action = enable ? '启用' : '禁用';
    
    fetch(`${API_BASE}/tasks/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ Enabled: enable })
    })
        .then(response => {
            if (response.ok) {
                loadTasks();
                alert(`任务${action}成功`);
            } else {
                throw new Error(`${action}失败`);
            }
        })
        .catch(error => {
            console.error(`${action}任务失败:`, error);
            alert(`${action}任务失败`);
        });
}

// 立即执行任务
async function executeTask(id) {
    try {
        const response = await fetch(`${API_BASE}/tasks/${id}/execute`, {
            method: 'POST'
        });
        
        if (response.ok) {
            alert('任务执行成功');
            loadTasks(); // 重新加载任务列表
        } else {
            const error = await response.json();
            alert('任务执行失败: ' + error.error);
        }
    } catch (error) {
        console.error('执行任务失败:', error);
        alert('执行任务失败: ' + error.message);
    }
}

// 显示任务执行日志
async function showTaskLogs(taskId) {
    try {
        const response = await fetch(`${API_BASE}/tasks/${taskId}/logs`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const logs = await response.json();
        
        // 创建模态框显示日志
        let logsHtml = '<h3>任务执行日志</h3>';
        if (logs.length === 0) {
            logsHtml += '<p>暂无日志</p>';
        } else {
            logsHtml += `
                <table>
                    <thead>
                        <tr>
                            <th>开始时间</th>
                            <th>结束时间</th>
                            <th>状态</th>
                            <th>输出</th>
                            <th>错误</th>
                        </tr>
                    </thead>
                    <tbody>
            `;
            
            logs.forEach(log => {
                logsHtml += `
                    <tr>
                        <td>${new Date(log.StartTime).toLocaleString()}</td>
                        <td>${new Date(log.EndTime).toLocaleString()}</td>
                        <td>${log.Success ? '成功' : '失败'}</td>
                        <td>${log.Output || ''}</td>
                        <td>${log.Error || ''}</td>
                    </tr>
                `;
            });
            
            logsHtml += '</tbody></table>';
        }
        
        // 显示模态框
        document.getElementById('modal-title').textContent = '任务日志';
        document.getElementById('modal-body').innerHTML = logsHtml;
        document.getElementById('modal').style.display = 'block';
    } catch (error) {
        console.error('获取任务日志失败:', error);
        alert('获取任务日志失败: ' + error.message);
    }
}

// 处理任务表单提交
function handleTaskSubmit(event) {
    event.preventDefault();
    
    const task = {
        Name: document.getElementById('taskName').value,
        Schedule: document.getElementById('taskSchedule').value,
        Command: document.getElementById('taskCommand').value,
        Method: document.getElementById('taskMethod').value,
        Headers: document.getElementById('taskHeaders').value,
        Description: document.getElementById('taskDescription').value,
        Enabled: document.getElementById('taskEnabled').checked
    };
    
    const method = editingTaskId ? 'PUT' : 'POST';
    const url = editingTaskId ? `${API_BASE}/tasks/${editingTaskId}` : `${API_BASE}/tasks`;
    
    fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(task)
    })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error('保存失败');
            }
        })
        .then(() => {
            loadTasks();
            resetForm();
            alert('任务保存成功');
        })
        .catch(error => {
            console.error('保存任务失败:', error);
            alert('保存任务失败');
        });
}

// 取消编辑
function cancelEdit() {
    resetForm();
}

// 重置表单
function resetForm() {
    document.getElementById('taskForm').reset();
    document.getElementById('taskId').value = '';
    editingTaskId = null;
    
    // 隐藏HTTP配置
    document.getElementById('httpConfig').style.display = 'none';
}

// 切换HTTP配置显示
function toggleHttpConfig() {
    const command = document.getElementById('taskCommand').value;
    const httpConfig = document.getElementById('httpConfig');
    const showHttpConfig = document.getElementById('showHttpConfig').value;
    
    // 根据下拉框选项决定是否显示HTTP配置
    if (showHttpConfig === 'show') {
        httpConfig.style.display = 'block';
    } else if (showHttpConfig === 'hide') {
        httpConfig.style.display = 'none';
    } else { // auto
        // 简单检查是否为URL
        if (command.startsWith('http://') || command.startsWith('https://')) {
            httpConfig.style.display = 'block';
        } else {
            httpConfig.style.display = 'none';
        }
    }
}

// 手动切换HTTP配置显示
function toggleHttpConfigManually() {
    toggleHttpConfig();
}

// 在页面加载完成后添加命令输入事件监听器
document.getElementById('taskCommand').addEventListener('input', function() {
    // 根据命令是否为URL显示HTTP配置
    toggleHttpConfig();
});