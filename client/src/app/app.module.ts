import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';

// used to create fake backend
import { fakeBackendProvider } from './_helpers';
import { AppComponent } from './app.component';
import { appRoutingModule } from './app.routing';
import { BasicAuthInterceptor, ErrorInterceptor } from './_helpers';
import { HomeComponent } from './home';
import { LoginComponent } from './login';
import { AccountsComponent } from './accounts/accounts.component';
import { ClassesComponent } from './classes/classes.component';
import { AccountDetailComponent } from './account-detail/account-detail.component';
import { RegisterComponent } from './register/register.component';
import { NewclassComponent } from './newclass/newclass.component';
import { ClassDetailComponent } from './class-detail/class-detail.component';
import { ExerciseDetailComponent } from './exercise-detail/exercise-detail.component';
import { FileUploadComponent } from './app-file-upload/app-file-upload.component';;
import { ChatComponent } from './chat/chat.component'

@NgModule({
    imports: [
        BrowserModule,
        ReactiveFormsModule,
        HttpClientModule,
        appRoutingModule,
        FormsModule
    ],
    declarations: [
        AppComponent,
        HomeComponent,
        LoginComponent
,
        AccountsComponent ,
        ClassesComponent ,
        AccountDetailComponent ,
        RegisterComponent,
        NewclassComponent
,
        ClassDetailComponent ,
        ExerciseDetailComponent ,
        FileUploadComponent ,
        ChatComponent     ],
    providers: [
        { provide: HTTP_INTERCEPTORS, useClass: BasicAuthInterceptor, multi: true },
        { provide: HTTP_INTERCEPTORS, useClass: ErrorInterceptor, multi: true },

        // provider used to create fake backend
        //fakeBackendProvider
    ],
    bootstrap: [AppComponent]
})
export class AppModule { }
