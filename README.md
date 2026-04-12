![Image not loading....](https://github.com/whatthexampp/webonly/blob/main/branding.png?raw=true)

Documentation soon  

How to build:  
go build -o daemon.exe src\Main.go (RUN THIS!!!)  

compile the Apache Module  
cl.exe /nologo /I"C:\(DIRECTORY TO ROOT OF APACHE)\include" /I"C:\(DIRECTORY TO ROOT OF APACHE)\include\apr" /MD /LD src\ModWebonly.c /Fe:mod_webonly.so /link /LIBPATH:"C:\(DIRECTORY TO ROOT OF APACHE)\lib" libhttpd.lib libapr-1.lib libaprutil-1.lib ws2_32.lib  

e.g. if you are using xampp:  
cl.exe /nologo /I"C:\xampp\apache\include" /I"C:\xampp\apache\include\apr" /MD /LD src\ModWebonly.c /Fe:mod_webonly.so /link /LIBPATH:"C:\xampp\apache\lib" libhttpd.lib libapr-1.lib libaprutil-1.lib ws2_32.lib  

copy the compiled module (mod_webonly.so) to C:\(DIRECTORY TO ROOT OF APACHE)\modules\  

e.g. xampp:  
C:\xampp\apache\modules\[...]  

then go to C:\(DIRECTORY TO ROOT OF APACHE)\conf and open httpd.conf  

search for "LoadModule" and at the end of the block add "LoadModule WebonlyModule modules/mod_webonly.so"  

scroll down all the way at the end of the file and add:  

<FilesMatch "\.wo$">
    SetHandler webonly-script
</FilesMatch>

and now it's done, run the daemon first so it actually works  
a
