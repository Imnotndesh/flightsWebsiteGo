document.addEventListener('DOMContentLoaded', function (){
    const form = document.getElementById('loginForm');
    form.addEventListener('submit', async function(event){
        event.preventDefault();
        const formData = new FormData(form);
        const data = {
            username: formData.get('Username'),
            password: formData.get('Password')
        };
        try{
            const response = await fetch('/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });
            const result = await response.text();
            console.log(result);
            if (result.status === 200){
                console.log("Success");
                window.location.href="/flights";
            }else{
                console.log("Failed to log in");
                location.reload();
            }
        }catch (error){
            console.log("Failed to log in or server side issues",error);
        }
    });
});