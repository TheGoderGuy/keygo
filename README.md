# KeyGo 🔑🚀  
_A simple tool for creating Keycloak users via client credentials._

## 📌 Prerequisites  
Before using `keygo`, **configure Keycloak properly**:

### 1️⃣ Enable Service Accounts  
1. **Log in to Keycloak Admin Panel**:  
   📌 [http://localhost:8080/admin](http://localhost:8080/admin)  
2. **Go to**: `Clients → user-manager`  
3. **Click on** the `"Settings"` tab  
4. **Set the following options**:  
   - ✅ **Access Type** → `confidential`  
   - ✅ **Service Accounts Enabled** → `ON`  
5. **Click Save** ✅  

---

### 2️⃣ Assign `manage-users` Role to the Client  
1. **Go to**: `Clients → user-manager → Service Account Roles`  
2. **Assign the correct role**:  
   - Click **"Assign Role"**  
   - Select **`realm-management/manage-users`**  
3. **Click Save** ✅  

---

## 📌 Set Environment Variables  
Instead of passing parameters every time, set them as environment variables:  

## 🚀 How to Use  
With environment variables set:  
```sh
go run main.go newuser user@example.com mypassword  
```
Or pass parameters via CLI:  

```sh 
go run main.go --keycloak-url="http://my-keycloak.com" \  
               --realm="myrealm" \  
               --client-id="user-manager" \  
               --client-secret="my-secret" \  
               newuser user@example.com mypassword  
```
---

## 📌 Running as an Executable  
1. Build the binary:  

go build -o createuser main.go  

2. Run the executable:  
```sh
./createuser newuser user@example.com newpassword  
```
---

## 📌 Expected Output  
User created successfully with an initial password (must be changed on first login)!  

Now you have a fully configured, easy-to-use Keycloak user creation tool! 🚀  

---

## 📌 Alternative: `.envrc` File for Auto-Loading Env Vars  
You can also create a `.envrc` file to automatically load environment variables when you enter your project directory:  

1. Create a file named `.envrc` and add:  
```sh
export KEYCLOAK_URL="http://localhost:8080"  
export KEYCLOAK_REALM="myrealm"  
export KEYCLOAK_CLIENT_ID="user-manager"  
export KEYCLOAK_CLIENT_SECRET="your-client-secret"  
```
2. Run:  

```sh
direnv allow .  
```

3. Now, every time you `cd` into the directory, your environment variables will be loaded automatically! ✅  

---

🔨 Happy coding! 🤖 Let me know if you need improvements. 🎉   
