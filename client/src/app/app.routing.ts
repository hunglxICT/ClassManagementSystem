import { Routes, RouterModule } from '@angular/router';

import { HomeComponent } from './home';
import { LoginComponent } from './login';
import { AuthGuard } from './_helpers';
import { AccountDetailComponent } from './account-detail/account-detail.component';
import { ClassDetailComponent } from './class-detail/class-detail.component';
import { RegisterComponent } from './register/register.component';
import { ClassesComponent } from './classes/classes.component';
import { NewclassComponent } from './newclass/newclass.component';
import { ExerciseDetailComponent } from './exercise-detail/exercise-detail.component';

const routes: Routes = [
    { path: '', component: HomeComponent, canActivate: [AuthGuard] },
    { path: 'login', component: LoginComponent },
    { path: 'detail/:id', component: AccountDetailComponent },
    { path: 'register', component: RegisterComponent },
    { path: 'newclass', component: NewclassComponent },
    { path: 'classdetail/:id', component: ClassDetailComponent },
    { path: 'exercise/:id', component: ExerciseDetailComponent },

    // otherwise redirect to home
    { path: '**', redirectTo: '' }
];

export const appRoutingModule = RouterModule.forRoot(routes);
