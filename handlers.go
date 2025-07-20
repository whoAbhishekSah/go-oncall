package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>OnCall Scheduler</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            color: #333;
        }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        
        .header {
            text-align: center;
            color: white;
            margin-bottom: 40px;
            padding: 40px 0;
        }
        .header h1 { font-size: 3rem; margin-bottom: 10px; }
        .header p { font-size: 1.2rem; opacity: 0.9; }
        
        .nav-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }
        
        .nav-card {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 12px;
            padding: 25px;
            text-decoration: none;
            color: #333;
            transition: all 0.3s ease;
            border: 2px solid transparent;
            backdrop-filter: blur(10px);
        }
        
        .nav-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
            border-color: #667eea;
            color: #333;
            text-decoration: none;
        }
        
        .nav-card h3 {
            font-size: 1.4rem;
            margin-bottom: 10px;
            color: #667eea;
        }
        
        .nav-card p {
            color: #666;
            line-height: 1.5;
        }
        
        .section {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 12px;
            padding: 30px;
            margin-bottom: 30px;
            backdrop-filter: blur(10px);
            display: none;
        }
        
        .section.active { display: block; }
        
        .form-group { margin-bottom: 20px; }
        
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #555;
        }
        
        input, textarea, select {
            width: 100%;
            padding: 12px 16px;
            border: 2px solid #e1e5e9;
            border-radius: 8px;
            font-size: 14px;
            transition: border-color 0.3s ease;
            background: white;
        }
        
        input:focus, textarea:focus, select:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }
        
        button {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 12px 24px;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        
        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        
        .back-btn {
            background: #6c757d;
            margin-bottom: 20px;
            padding: 8px 16px;
            font-size: 14px;
        }
        
        .back-btn:hover {
            background: #5a6268;
        }
        
        .item-card {
            background: #f8f9fa;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 15px;
            border-left: 4px solid #667eea;
        }
        
        .item-card h4 {
            color: #667eea;
            margin-bottom: 10px;
        }
        
        .item-card p {
            margin-bottom: 5px;
            color: #666;
        }
        
        .emoji { font-size: 1.5rem; margin-right: 10px; }
        
        /* Toast notification styles */
        .toast-container {
            position: fixed;
            top: 20px;
            right: 20px;
            z-index: 1000;
            display: flex;
            flex-direction: column;
            gap: 10px;
        }
        
        .toast {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 8px;
            padding: 16px 20px;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
            backdrop-filter: blur(10px);
            border-left: 4px solid;
            min-width: 300px;
            transform: translateX(100%);
            opacity: 0;
            transition: all 0.3s ease;
        }
        
        .toast.show {
            transform: translateX(0);
            opacity: 1;
        }
        
        .toast.success {
            border-left-color: #28a745;
        }
        
        .toast.error {
            border-left-color: #dc3545;
        }
        
        .toast-header {
            font-weight: 600;
            margin-bottom: 4px;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .toast-message {
            color: #666;
            font-size: 14px;
        }
        
        .toast-close {
            position: absolute;
            top: 8px;
            right: 12px;
            background: none;
            border: none;
            font-size: 18px;
            cursor: pointer;
            color: #999;
            padding: 0;
            width: 20px;
            height: 20px;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        
        .toast-close:hover {
            color: #666;
        }
    </style>
</head>
<body>
    <!-- Toast Container -->
    <div class="toast-container" id="toast-container"></div>
    
    <div class="container">
        <div class="header">
            <h1>🚨 OnCall Scheduler</h1>
            <p>Manage your team's on-call rotations with ease</p>
        </div>

        <!-- Navigation Grid -->
        <div id="nav-section">
            <div class="nav-grid">
                <a href="#" class="nav-card" onclick="showSection('create-user')">
                    <h3><span class="emoji">👤</span>Create User</h3>
                    <p>Add new team members with email and Slack details</p>
                </a>
                
                <a href="#" class="nav-card" onclick="showSection('create-team')">
                    <h3><span class="emoji">👥</span>Create Team</h3>
                    <p>Set up new teams for organizing your on-call rotations</p>
                </a>
                
                <a href="#" class="nav-card" onclick="showSection('create-schedule')">
                    <h3><span class="emoji">📅</span>Create Schedule</h3>
                    <p>Design on-call schedules with rotation periods</p>
                </a>
                
                <a href="#" class="nav-card" onclick="showSection('view-users')">
                    <h3><span class="emoji">📋</span>View Users</h3>
                    <p>Browse all registered users and their details</p>
                </a>
                
                <a href="#" class="nav-card" onclick="showSection('view-teams')">
                    <h3><span class="emoji">🏢</span>View Teams</h3>
                    <p>See all teams and their members</p>
                </a>
                
                <a href="#" class="nav-card" onclick="showSection('view-schedules')">
                    <h3><span class="emoji">⏰</span>View Schedules</h3>
                    <p>Monitor active schedules and rotations</p>
                </a>
            </div>
        </div>

        <!-- Create User Section -->
        <div id="create-user" class="section">
            <button class="back-btn" onclick="showNav()">← Back to Menu</button>
            <h2>👤 Create User</h2>
            <form id="userForm">
                <div class="form-group">
                    <label for="userEmail">Email Address:</label>
                    <input type="email" id="userEmail" name="userEmail" placeholder="john.doe@company.com" required>
                </div>
                <div class="form-group">
                    <label for="userSlackHandle">Slack Handle:</label>
                    <input type="text" id="userSlackHandle" name="userSlackHandle" placeholder="@johndoe" required>
                </div>
                <div class="form-group">
                    <label for="userTeamId">Team ID:</label>
                    <input type="number" id="userTeamId" name="userTeamId" placeholder="1" required>
                </div>
                <button type="submit">Create User</button>
            </form>
        </div>

        <!-- Create Team Section -->
        <div id="create-team" class="section">
            <button class="back-btn" onclick="showNav()">← Back to Menu</button>
            <h2>👥 Create Team</h2>
            <form id="teamForm">
                <div class="form-group">
                    <label for="teamName">Team Name:</label>
                    <input type="text" id="teamName" name="teamName" placeholder="Backend Engineering" required>
                </div>
                <button type="submit">Create Team</button>
            </form>
        </div>

        <!-- Create Schedule Section -->
        <div id="create-schedule" class="section">
            <button class="back-btn" onclick="showNav()">← Back to Menu</button>
            <h2>📅 Create Schedule</h2>
            <form id="scheduleForm">
                <div class="form-group">
                    <label for="scheduleName">Schedule Name:</label>
                    <input type="text" id="scheduleName" name="scheduleName" placeholder="Weekend On-Call Rotation" required>
                </div>
                <div class="form-group">
                    <label for="teamId">Team ID:</label>
                    <input type="number" id="teamId" name="teamId" placeholder="1" required>
                </div>
                <div class="form-group">
                    <label for="startTime">Start Time:</label>
                    <input type="datetime-local" id="startTime" name="startTime" required>
                </div>
                <div class="form-group">
                    <label for="endTime">End Time:</label>
                    <input type="datetime-local" id="endTime" name="endTime" required>
                </div>
                <div class="form-group">
                    <label for="rotationPeriod">Rotation Period:</label>
                    <div style="display: flex; gap: 10px; align-items: center;">
                        <div style="flex: 1;">
                            <label for="rotationDays" style="font-size: 12px; margin-bottom: 2px;">Days</label>
                            <input type="number" id="rotationDays" name="rotationDays" min="0" max="365" placeholder="0" style="width: 100%;">
                        </div>
                        <div style="flex: 1;">
                            <label for="rotationHours" style="font-size: 12px; margin-bottom: 2px;">Hours</label>
                            <input type="number" id="rotationHours" name="rotationHours" min="0" max="23" placeholder="1" style="width: 100%;">
                        </div>
                        <div style="flex: 1;">
                            <label for="rotationMinutes" style="font-size: 12px; margin-bottom: 2px;">Minutes</label>
                            <input type="number" id="rotationMinutes" name="rotationMinutes" min="0" max="59" placeholder="0" style="width: 100%;">
                        </div>
                        <div style="flex: 1;">
                            <label for="rotationSeconds" style="font-size: 12px; margin-bottom: 2px;">Seconds</label>
                            <input type="number" id="rotationSeconds" name="rotationSeconds" min="0" max="59" placeholder="0" style="width: 100%;">
                        </div>
                    </div>
                    <small style="color: #666; font-size: 12px; margin-top: 5px; display: block;">
                        Minimum rotation period is 1 second. Common examples: 1 day = 1d 0h 0m 0s, 8 hours = 0d 8h 0m 0s
                    </small>
                </div>
                <div class="form-group">
                    <label for="participants">Participants (comma-separated user IDs):</label>
                    <textarea id="participants" name="participants" rows="3" placeholder="1,2,3" required></textarea>
                </div>
                <button type="submit">Create Schedule</button>
            </form>
        </div>

        <!-- View Users Section -->
        <div id="view-users" class="section">
            <button class="back-btn" onclick="showNav()">← Back to Menu</button>
            <h2>📋 All Users</h2>
            <div id="usersList"></div>
        </div>

        <!-- View Teams Section -->
        <div id="view-teams" class="section">
            <button class="back-btn" onclick="showNav()">← Back to Menu</button>
            <h2>🏢 All Teams</h2>
            <div id="teamsList"></div>
        </div>

        <!-- View Schedules Section -->
        <div id="view-schedules" class="section">
            <button class="back-btn" onclick="showNav()">← Back to Menu</button>
            <h2>⏰ All Schedules</h2>
            <div id="schedulesList"></div>
        </div>
    </div>

    <script>
        // Toast notification system
        function showToast(type, title, message, duration = 4000) {
            const container = document.getElementById('toast-container');
            const toast = document.createElement('div');
            toast.className = 'toast ' + type;
            
            toast.innerHTML = 
                '<button class="toast-close" onclick="removeToast(this.parentElement)">&times;</button>' +
                '<div class="toast-header">' +
                    (type === 'success' ? '✅' : '❌') + ' ' + title +
                '</div>' +
                '<div class="toast-message">' + message + '</div>';
            
            container.appendChild(toast);
            
            // Trigger animation
            setTimeout(function() { toast.classList.add('show'); }, 100);
            
            // Auto remove
            setTimeout(function() { removeToast(toast); }, duration);
        }
        
        function removeToast(toast) {
            toast.classList.remove('show');
            setTimeout(function() {
                if (toast.parentElement) {
                    toast.parentElement.removeChild(toast);
                }
            }, 300);
        }
        
        // Utility function to format duration from seconds
        function formatDuration(totalSeconds) {
            const days = Math.floor(totalSeconds / (24 * 60 * 60));
            const hours = Math.floor((totalSeconds % (24 * 60 * 60)) / (60 * 60));
            const minutes = Math.floor((totalSeconds % (60 * 60)) / 60);
            const seconds = totalSeconds % 60;
            
            const parts = [];
            if (days > 0) parts.push(days + 'd');
            if (hours > 0) parts.push(hours + 'h');
            if (minutes > 0) parts.push(minutes + 'm');
            if (seconds > 0 || parts.length === 0) parts.push(seconds + 's');
            
            return parts.join(' ');
        }
        
        // Navigation functions
        function showSection(sectionId) {
            // Hide navigation
            document.getElementById('nav-section').style.display = 'none';
            
            // Hide all sections
            const sections = document.querySelectorAll('.section');
            sections.forEach(section => section.classList.remove('active'));
            
            // Show selected section
            const targetSection = document.getElementById(sectionId);
            if (targetSection) {
                targetSection.classList.add('active');
                
                // Load data for view sections
                if (sectionId === 'view-users') loadUsers();
                if (sectionId === 'view-teams') loadTeams();
                if (sectionId === 'view-schedules') loadSchedules();
            }
        }
        
        function showNav() {
            // Show navigation
            document.getElementById('nav-section').style.display = 'block';
            
            // Hide all sections
            const sections = document.querySelectorAll('.section');
            sections.forEach(section => section.classList.remove('active'));
        }
        
        // Form handlers
        document.getElementById('userForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const formData = new FormData(this);
            
            fetch('/users', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    email: formData.get('userEmail'),
                    slack_handle: formData.get('userSlackHandle'),
                    team_id: parseInt(formData.get('userTeamId'))
                })
            })
            .then(response => response.json())
            .then(data => {
                showToast('success', 'User Created', 'User has been successfully added to the system!');
                this.reset();
            })
            .catch(error => {
                console.error('Error:', error);
                showToast('error', 'Creation Failed', 'Failed to create user. Please try again.');
            });
        });

        document.getElementById('teamForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const formData = new FormData(this);
            
            fetch('/teams', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    name: formData.get('teamName')
                })
            })
            .then(response => response.json())
            .then(data => {
                showToast('success', 'Team Created', 'Team has been successfully created!');
                this.reset();
            })
            .catch(error => {
                console.error('Error:', error);
                showToast('error', 'Creation Failed', 'Failed to create team. Please try again.');
            });
        });

        document.getElementById('scheduleForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const formData = new FormData(this);
            const participants = formData.get('participants').split(',').map(p => parseInt(p.trim()));
            
            // Calculate total rotation period in seconds
            const days = parseInt(formData.get('rotationDays') || 0);
            const hours = parseInt(formData.get('rotationHours') || 0);
            const minutes = parseInt(formData.get('rotationMinutes') || 0);
            const seconds = parseInt(formData.get('rotationSeconds') || 0);
            
            const totalSeconds = (days * 24 * 60 * 60) + (hours * 60 * 60) + (minutes * 60) + seconds;
            
            if (totalSeconds < 1) {
                showToast('error', 'Invalid Rotation Period', 'Rotation period must be at least 1 second.');
                return;
            }
            
            fetch('/schedules', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    team_id: parseInt(formData.get('teamId')),
                    name: formData.get('scheduleName'),
                    start_time: formData.get('startTime'),
                    end_time: formData.get('endTime'),
                    rotation_period: totalSeconds,
                    participants: participants
                })
            })
            .then(response => response.json())
            .then(data => {
                showToast('success', 'Schedule Created', 'On-call schedule has been successfully created!');
                this.reset();
            })
            .catch(error => {
                console.error('Error:', error);
                showToast('error', 'Creation Failed', 'Failed to create schedule. Please check your inputs and try again.');
            });
        });

        function loadUsers() {
            fetch('/users')
                .then(response => response.json())
                .then(users => {
                    const usersList = document.getElementById('usersList');
                    if (users.length === 0) {
                        usersList.innerHTML = '<p>No users found. <a href="#" onclick="showSection(\'create-user\')">Create your first user</a></p>';
                    } else {
                        usersList.innerHTML = users.map(user => 
                            '<div class="item-card">' +
                            '<h4>👤 ' + user.email + ' (ID: ' + user.id + ')</h4>' +
                            '<p><strong>Slack:</strong> ' + user.slack_handle + '</p>' +
                            '<p><strong>Team ID:</strong> ' + user.team_id + '</p>' +
                            '<p><strong>Created:</strong> ' + new Date(user.created_at).toLocaleString() + '</p>' +
                            '</div>'
                        ).join('');
                    }
                });
        }

        function loadTeams() {
            fetch('/teams')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('HTTP ' + response.status + ': ' + response.statusText);
                    }
                    return response.json();
                })
                .then(teams => {
                    console.log('Teams loaded:', teams);
                    const teamsList = document.getElementById('teamsList');
                    if (!teams || teams.length === 0) {
                        teamsList.innerHTML = '<p>No teams found. <a href="#" onclick="showSection(\'create-team\')">Create your first team</a></p>';
                    } else {
                        teamsList.innerHTML = teams.map(team => 
                            '<div class="item-card">' +
                            '<h4>👥 ' + team.name + ' (ID: ' + team.id + ')</h4>' +
                            '<p><strong>Members:</strong> ' + (team.users && team.users.length > 0 ? team.users.map(u => u.email + ' (' + u.slack_handle + ')').join(', ') : 'No members yet') + '</p>' +
                            '<p><strong>Created:</strong> ' + new Date(team.created_at).toLocaleString() + '</p>' +
                            '</div>'
                        ).join('');
                    }
                })
                .catch(error => {
                    console.error('Error loading teams:', error);
                    const teamsList = document.getElementById('teamsList');
                    teamsList.innerHTML = '<p>Error loading teams: ' + error.message + '</p>';
                });
        }

        function loadSchedules() {
            fetch('/schedules')
                .then(response => response.json())
                .then(schedules => {
                    const schedulesList = document.getElementById('schedulesList');
                    if (schedules.length === 0) {
                        schedulesList.innerHTML = '<p>No schedules found. <a href="#" onclick="showSection(\'create-schedule\')">Create your first schedule</a></p>';
                    } else {
                        schedulesList.innerHTML = schedules.map(schedule => 
                            '<div class="item-card">' +
                            '<h4>📅 ' + schedule.name + ' (Team ID: ' + schedule.team_id + ')</h4>' +
                            '<p><strong>Start:</strong> ' + new Date(schedule.start_time).toLocaleString() + '</p>' +
                            '<p><strong>End:</strong> ' + new Date(schedule.end_time).toLocaleString() + '</p>' +
                            '<p><strong>Rotation:</strong> ' + formatDuration(schedule.rotation_period) + '</p>' +
                            '<p><strong>Participants (User IDs):</strong> ' + (schedule.participants ? schedule.participants.join(', ') : 'None') + '</p>' +
                            '</div>'
                        ).join('');
                    }
                });
        }

        // Initialize the page
        showNav();
    </script>
</body>
</html>
	`
	
	t := template.Must(template.New("home").Parse(tmpl))
	t.Execute(w, nil)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Email       string `json:"email"`
		SlackHandle string `json:"slack_handle"`
		TeamID      int    `json:"team_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	id, err := createUser(user.Email, user.SlackHandle, user.TeamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"id":      id,
		"message": "User created successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := getUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createTeamHandler(w http.ResponseWriter, r *http.Request) {
	var team struct {
		Name string `json:"name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	id, err := createTeam(team.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"id":      id,
		"message": "Team created successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getTeamsHandler(w http.ResponseWriter, r *http.Request) {
	teams, err := getTeams()
	if err != nil {
		log.Printf("Error getting teams: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	log.Printf("Retrieved %d teams", len(teams))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func createScheduleHandler(w http.ResponseWriter, r *http.Request) {
	var schedule struct {
		TeamID         int      `json:"team_id"`
		Name           string   `json:"name"`
		StartTime      string   `json:"start_time"`
		EndTime        string   `json:"end_time"`
		RotationPeriod int      `json:"rotation_period"`
		Participants   []int    `json:"participants"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	startTime, err := time.Parse("2006-01-02T15:04", schedule.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}
	
	endTime, err := time.Parse("2006-01-02T15:04", schedule.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}
	
	id, err := createSchedule(schedule.TeamID, schedule.Name, startTime, endTime, schedule.RotationPeriod, schedule.Participants)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"id":      id,
		"message": "Schedule created successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	schedules, err := getSchedules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}